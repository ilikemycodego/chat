package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitDB() {
	if DB != nil {
		return
	}

	// Пути к .env только для Linux
	paths := []string{
		"/var/www/auth/.env", // VPS
		"../.env",            // локально рядом с проектом / бинарником
	}

	envLoaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Println("📄 Загружен .env:", path)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		log.Println("⚠️ .env не найден — используются переменные среды")
	}

	// Получение данных из переменных среды
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	sslmode := os.Getenv("DB_SSLMODE")

	// Формируем DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		user, password, host, dbname, sslmode,
	)

	// Подключение к базе
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

	// Выполнение миграций (если есть файл)
	sqlFile := "db/migrations/auth.sql"
	if _, err := os.Stat(sqlFile); err == nil {
		sqlBytes, err := os.ReadFile(sqlFile)
		if err != nil {
			log.Fatalf("❌ Не удалось прочитать файл миграции: %v", err)
		}
		if _, err = DB.Exec(ctx, string(sqlBytes)); err != nil {
			log.Fatalf("❌ Ошибка миграции: %v", err)
		}
		log.Println("✅ Таблицы готовы!")
	} else {
		log.Println("⚠️ Файл миграций не найден:", sqlFile)
	}
}
