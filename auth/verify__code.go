package auth

import (
	"chat/db"
	"chat/token"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Проверка одноразового кода
func verifyCode(email, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int

	err := db.DB.QueryRow(ctx, `
		UPDATE auth_codes
		SET used = true
		WHERE email = $1
		  AND code = $2
		  AND used = false
		  AND created_at > NOW() - INTERVAL '10 minutes'
		RETURNING id
	`, email, code).Scan(&id)

	if err != nil {
		// код не найден / уже использован / истёк
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		// реальная ошибка БД
		return false, err
	}

	return true, nil
}

// Создание сессии на 7 дней
func createSession(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionID := uuid.New().String()
	expires := time.Now().Add(7 * 24 * time.Hour)

	_, err := db.DB.Exec(ctx, `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, sessionID, userID, expires)

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Поиск или создание нового пользователя (без роли)
func getOrCreateUser(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string

	// пробуем найти
	err := db.DB.QueryRow(ctx,
		`SELECT id FROM users WHERE email = $1`,
		email,
	).Scan(&id)

	if err == nil {
		return id, nil
	}

	// создаём (без падения при гонке)
	newID := uuid.New().String()

	_, err = db.DB.Exec(ctx, `
		INSERT INTO users (id, email, name)
		VALUES ($1, $2, 'amigo')
		ON CONFLICT (email) DO NOTHING
	`, newID, email)
	if err != nil {
		return "", err
	}

	// гарантированно получаем id
	err = db.DB.QueryRow(ctx,
		`SELECT id FROM users WHERE email = $1`,
		email,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

// Верификация кода + создание пользователя + создание сессии

func VerifyCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.FormValue("email")
		code := r.FormValue("code")

		if email == "" || code == "" {
			fmt.Fprint(w, "Заполните все поля!")
			return
		}

		ok, err := verifyCode(email, code)
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка проверки кода:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		if !ok {
			fmt.Fprint(w, "Неверный или уже использованный код!")
			return
		}

		// 1️⃣ Создаём или ищем пользователя
		userID, err := getOrCreateUser(email)
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка getOrCreateUser:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		// 2️⃣ Создаём сессию
		sessionID, err := createSession(userID)
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка createSession:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		// 3️⃣ Генерируем JWT только с sessionID
		jwtToken, err := token.GenerateJWT(sessionID, token.ExpirationTime())
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка GenerateJWT:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		// 4️⃣ Устанавливаем cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "jwtAuth",
			Value:    jwtToken,
			HttpOnly: true,
			Secure:   true, // если HTTPS
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			MaxAge:   7 * 24 * 3600,
		})

		log.Println("[VerifyCodeHandler] JWT создан и cookie установлено")

		// 5️⃣ HTMX редирект
		w.Header().Set("HX-Redirect", "/")
		fmt.Fprint(w, "OK")
	}
}
