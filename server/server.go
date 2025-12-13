package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NewServer собирает все шаблоны и маршруты
func NewServer() http.Handler {
	// Загружаем все шаблоны из templates/**/*
	tmpl := template.Must(template.ParseGlob("templates/**/*.html"))

	m := mux.NewRouter()

	// Обслуживание статики
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Регистрируем все маршруты, включая /theme
	RegisterRoutes(m, tmpl)

	log.Println("✅ Сервер готов!")
	return m
}
