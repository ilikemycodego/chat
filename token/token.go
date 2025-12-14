package token

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	jwtSecret []byte
	once      sync.Once
)

// loadSecret — загружает JWT_SECRET из .env
func loadSecret() {
	paths := []string{
		"/var/www/auth/.env", // VPS
		"../.env",            // локально рядом с проектом / бинарником
	}

	envLoaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		panic("❌ JWT_SECRET не найден в .env или окружении")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("❌ JWT_SECRET не найден в .env или окружении")
	}
	jwtSecret = []byte(secret)
}

// ExpirationTime — единый источник времени жизни токена/сессии
func ExpirationTime() time.Time {
	return time.Now().Add(7 * 24 * time.Hour) // 7 дней
}

// Claims — структура данных, хранимая внутри токена
type Claims struct {
	SessionID string `json:"session_id"`

	jwt.RegisteredClaims
}

// GenerateJWT — создаёт новый токен по sessionID
func GenerateJWT(sessionID string, expiration time.Time) (string, error) {
	once.Do(loadSecret)

	claims := &Claims{
		SessionID: sessionID, // добавляем сюда
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(jwtSecret)
}

// ParseJWT — разбирает токен и возвращает Claims
func ParseJWT(tokenStr string) (*Claims, error) {
	once.Do(loadSecret)

	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("неверный токен")
	}

	return claims, nil
}
