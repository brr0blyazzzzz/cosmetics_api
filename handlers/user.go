package handlers

import (
	"cosmetics/models"
	"cosmetics/repository"
	"encoding/json"
	"log"
	"net/http"
)

// создание структуры-зависимости для взаимодействия с репозиторием
type UserHandler struct {
	Repo *repository.UserRepository
}

// конструктор для создания экземпляра с внедренным репозиторием
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// обработка POST-запроса к маршруту входа
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	//получение тела запроса и декодирование в структуру(username пользователя и пароль)
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	//поиск по имени пользователя(аутентификация)
	user, err := h.Repo.GetUserByUsername(reqUser.UserName)
	//обработка серверных ошибок
	if err != nil {
		log.Printf("Ошибка при поиске пользователя: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	//обработка ошибки данных пользователя
	if user == nil {
		handleUnauthorized(w, "Неверное имя пользователя или пароль")
		return
	}
	//обработка при неверном пароле
	if !user.CheckPassword(reqUser.Password) {
		handleUnauthorized(w, "Неверное имя пользователя или пароль")
		return
	}
	//генерация токена
	tokenString, err := generateToken(user.UserName)
	//обработка ошибки при генерации токена
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}
	//установка заголовка для браузера
	w.Header().Set("Content-Type", "application/json")
	//отправка json
	json.NewEncoder(w).Encode(models.Response{
		Message: "Авторизация прошла успешно",
		Data: map[string]string{
			"token": tokenString,
		},
	})
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	//получение тела запроса и декодирование в структуру(username пользователя и пароль)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	//обработка и проверка на пустоту полей
	if user.UserName == "" || user.Password == "" {
		http.Error(w, "Имя пользователя и пароль обязательны", http.StatusBadRequest)
		return
	}
	//обработка хеширования пароля
	if err := user.SetPassword(user.Password); err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	//создание пользователя и обработка ошибки создания записи в бд
	if err := h.Repo.CreateUser(&user); err != nil {
		log.Printf("Ошибка создания пользователя в БД: %v", err)
		http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
		return
	}
	//установка заголовка
	w.Header().Set("Content-Type", "application/json")
	//отправка json-ответа
	json.NewEncoder(w).Encode(models.Response{
		Message: "Регистрация прошла успешно",
	})
}

// функция для обработки неавторизованных пользователей
func handleUnauthorized(w http.ResponseWriter, message string) {
	//установка заголовка
	w.Header().Set("Content-Type", "application/json")
	//отправка json-ответа с кодом 401
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error:   "Unauthorized",
		Message: message,
	})
}
