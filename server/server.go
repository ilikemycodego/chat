package server

import (
	"html/template"
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

// NewServer собирает все шаблоны, роуты и возвращает готовый mux
func NewServer() http.Handler {
	// Загружаем все шаблоны
	tmpl := template.Must(template.ParseGlob("templates/**/*.html"))

	// Инициализируем роутер
	m := mux.NewRouter()

	// Обслуживаем статику
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Подключаем маршруты
	RegisterRoutes(m, tmpl)

	log.Println("✅ Gorilla mux готов!")
	return m
}
