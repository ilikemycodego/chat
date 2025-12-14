package db

import (
	"context"
	"time"
)

// Сессия
type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

// Пользователь
type User struct {
	ID   string
	Name string
}

// Получаем сессию по ID
func GetSession(sessionID string) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := &Session{}
	err := DB.QueryRow(ctx, `
		SELECT id, user_id, expires_at
		FROM sessions
		WHERE id = $1
		  AND expires_at > NOW()
	`, sessionID).Scan(&s.ID, &s.UserID, &s.ExpiresAt)

	if err != nil {
		return nil, err
	}
	return s, nil
}

// Получаем пользователя по ID
func GetUserByID(userID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := &User{}
	err := DB.QueryRow(ctx, `
	SELECT id, name FROM users WHERE id=$1
    `, userID).Scan(&u.ID, &u.Name)

	if err != nil {
		return nil, err
	}

	return u, nil
}
