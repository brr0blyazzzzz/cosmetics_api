package handlers

import (
	"cosmetics/models"
	"cosmetics/repository"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// структура данных для рендеринга главной и админ-страниц
type WelcomePageData struct {
	Products               []models.Product
	Manufacturers          []models.Manufacturer
	SelectedManufacturerID int
	SearchQuery            string
	IsAuthenticated        bool
}

// обработчик главной страницы
func WebHandler(productRepo *repository.ProductRepository, manufacturerRepo *repository.ManufacturerRepository) http.HandlerFunc {
	// предварительная загрузка и парсинг шаблонов при старте приложения
	tmpl, err := template.ParseFiles("views/index.html", "views/admin.html")
	// обработка ошибки загрузки шаблонов
	if err != nil {
		log.Fatalf("Ошибка загрузки страницы: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// извлечение идентификатора производителя из параметров URL для фильтрации
		manufacturerIDStr := r.URL.Query().Get("manufacturer_id")
		// преобразование идентификатора производителя в int
		manufacturerID, _ := strconv.Atoi(manufacturerIDStr)
		// извлечение и очистка строки поиска из параметров URL
		searchQuery := strings.TrimSpace(r.URL.Query().Get("query"))

		// получение отфильтрованных продуктов из репозитория
		products, err := productRepo.GetProductsSearch(manufacturerID, searchQuery)
		// обработка ошибки получения данных о продуктах
		if err != nil {
			http.Error(w, "Ошибка получения данных о продуктах", http.StatusInternalServerError)
			return
		}

		// получение списка всех производителей
		manufacturers, err := manufacturerRepo.GetAll()
		// обработка ошибки получения производителей (не фатальная, просто логируется)
		if err != nil {
			log.Printf("Ошибка получения производителей: %v", err)
			// инициализация пустым списком, чтобы не сломать шаблон
			manufacturers = []models.Manufacturer{}
		}

		// Проверка авторизации по JWT из cookie
		isAuthenticated := false

		// попытка получить cookie "token"
		if cookie, err := r.Cookie("token"); err == nil {
			claims := &Claims{}
			// парсинг токена с проверкой подписи и извлечением клеймов
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			// проверка: если нет ошибки парсинга и токен валиден (срок действия, подпись)
			if err == nil && token.Valid {
				isAuthenticated = true
			}
		}

		// сбор всех данных в структуру для шаблона
		data := WelcomePageData{
			Products:               products,
			Manufacturers:          manufacturers,
			SelectedManufacturerID: manufacturerID,
			SearchQuery:            searchQuery,
			IsAuthenticated:        isAuthenticated, // флаг для условного рендеринга
		}

		// установка заголовка ответа
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// выполнение шаблона index с передачей данных
		if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
			// обработка ошибки выполнения шаблона
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}
}

// обработчик админ-панели
func AdminHandler(productRepo *repository.ProductRepository) http.HandlerFunc {
	// предварительная загрузка и парсинг шаблонов при старте приложения
	tmpl, err := template.ParseFiles("views/index.html", "views/admin.html")
	// обработка ошибки загрузки шаблонов
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// получение всех продуктов из репозитория для админ-панели
		products, err := productRepo.GetAll()
		// обработка ошибки получения данных о продуктах
		if err != nil {
			log.Printf("Ошибка получения данных о продуктах для админки: %v", err)
			http.Error(w, "Ошибка получения данных о продуктах", http.StatusInternalServerError)
			return
		}

		// создание структуры данных для шаблона
		// если пользователь попал на /admin, считаем его авторизованным
		data := WelcomePageData{
			Products:        products,
			IsAuthenticated: true, // устанавливаем в true, так как маршрут защищен Middleware
		}

		// установка заголовка ответа
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// выполнение шаблона "admin" с передачей данных
		if err := tmpl.ExecuteTemplate(w, "admin", data); err != nil {
			// обработка ошибки выполнения шаблона
			log.Printf("Ошибка выполнения шаблона 'admin': %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}
}

// Обработчик выхода пользователя
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// установка cookie token со сроком действия в прошлом
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "", // обнуление значения
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour), // устанавка истекшего времени
			HttpOnly: true,
		})
		// перенаправление на страницу входа
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
