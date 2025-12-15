package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitDB() {
	if DB != nil {
		return
	}

	// Пути к .env
	paths := []string{
		"/var/www/auth/.env",
		"../.env",
	}

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Println("📄 Загружен .env:", path)
			break
		}
	}

	// Переменные окружения
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		user, password, host, dbname, sslmode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	if err = DB.Ping(ctx); err != nil {
		log.Fatalf("❌ Нет соединения с БД: %v", err)
	}

	log.Println("✅ Подключение к БД установлено!")

	// ─────────────────────────────
	// Запуск миграций
	// ─────────────────────────────
	runMigrations(ctx)

	log.Println("✅ Таблицы готовы!")
}

func runMigrations(ctx context.Context) {
	dir := "db/migrations"

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("⚠️ Каталог миграций не найден:", dir)
		return
	}

	// Сортировка файлов: 001_, 002_ ...
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}

		path := filepath.Join(dir, e.Name())

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("❌ Ошибка чтения %s: %v", path, err)
		}

		if _, err := DB.Exec(ctx, string(sqlBytes)); err != nil {
			log.Fatalf("❌ Ошибка выполнения %s: %v", path, err)
		}

		log.Println("➡️ Миграция выполнена:", e.Name())
	}
}
