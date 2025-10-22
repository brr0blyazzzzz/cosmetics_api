package handlers

import (
	"context"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
)

//Выполняет поиск JWT в cookie запроса
//Если токен найден, идет проверка подлинности
//Если токена нет или он недействителен, редиект на login
//Если токен действителен, идет передача следующему обработчику

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Извлечение и проверка токена на существование
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		//Создание экземпляра структуры для хранения токена пользователя и извлечение jwt-ключа
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		//Проверка на валидность ключа
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		//Передача данных
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		//Вызов следующего обработчика
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
