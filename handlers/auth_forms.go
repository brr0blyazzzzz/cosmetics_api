package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

//Обработчики для HTTP запросов для регистрации и входа

// Отрисовка страницы входа
func (h *UserHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	//Загрузка и парсинг шаблона для логина
	tmpl, err := template.ParseFiles("views/login.html")
	//Обработка ошибок
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы входа", http.StatusInternalServerError)
		return
	}
	//Передача заголовка браузеру
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//Выполняет шаблон
	tmpl.ExecuteTemplate(w, "login", nil)
}

// Отрисовка страницы регистрации
func (h *UserHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	//Загрузка и парсинг шаблона для регистрации
	tmpl, err := template.ParseFiles("views/register.html")
	//Обработка ошибок
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы регистрации", http.StatusInternalServerError)
		return
	}
	//Передача заголовка браузеру
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//Выполняет шаблон
	tmpl.ExecuteTemplate(w, "register", nil)
}

// Обработка формы входа
func (h *UserHandler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	//Получение данных из формы
	username := r.FormValue("email")
	password := r.FormValue("password")
	//Поиск пользователя в базе данных(username приходит из AuthMiddleware)
	user, err := h.Repo.GetUserByUsername(username)
	//Обработка ошибок
	if err != nil {
		log.Printf("Ошибка при поиске пользователя: %v", err)
		http.Error(w, "Ошибка при поиске пользователя", http.StatusInternalServerError)
		return
	}
	//Сравнение пароля с хешем в базе данных и обработка ошибок
	if user == nil || !user.CheckPassword(password) {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	//Генерация JWT(здесь — user.UserName и криптографическая подпись)
	tokenString, err := generateToken(user.UserName)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	//Установка cookie для передачи данных
	http.SetCookie(w, &http.Cookie{
		Name:     "token",                        //имя, под которым будет храниться JWT
		Value:    tokenString,                    //сам JWT
		Path:     "/",                            //доступ к cookie всему сайту
		Expires:  time.Now().Add(24 * time.Hour), //срок действия ключа
		HttpOnly: true,                           //запрещает доступ клиенту
		Secure:   false,                          //HTTP-протокол
	})
	//Перенаправление на админ-панель
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
