package auth

import (
	"chat/db"
	"context"
	"fmt"
	"time"
)

func generateUniqueNick(base string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nick := sanitizeUsername(base)
	candidate := nick
	suffix := 1

	for {
		var exists bool
		err := db.DB.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM users WHERE nick = $1)`,
			candidate,
		).Scan(&exists)
		if err != nil {
			return "", err
		}

		if !exists {
			// нашли уникальный ник
			return candidate, nil
		}

		// добавляем число и проверяем снова
		candidate = fmt.Sprintf("%s%d", nick, suffix)
		suffix++
	}
}
