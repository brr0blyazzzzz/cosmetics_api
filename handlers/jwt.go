package handlers

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Ключ для криптографической подписи токенов
var jwtKey = []byte("sdcfghyjukilolikujyhgfdefrgthjkl;kjhgfdgthyjukilo;plkjyhtgrfghjkl;")

// Структура клеймов JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Создание JWT-токена для пользователя
func generateToken(username string) (string, error) {
	//Установка срока действия
	expirationTime := time.Now().Add(24 * time.Hour)
	//СОздание структуры клейма
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}
	//Создание токена с подписью
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}
