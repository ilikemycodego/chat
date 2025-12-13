package server

import (
	"chat/auth"
	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует все маршруты
func RegisterRoutes(m *mux.Router, tmpl *template.Template) {
	// Основной маршрут
	m.HandleFunc("/", auth.BaseHandler(tmpl))

	m.HandleFunc("/theme", auth.ToggleThemeHandler)

	// Пример дополнительных маршрутов:
	// r.HandleFunc("/login", auth.LoginHandler(tmpl))
	// r.HandleFunc("/logout", auth.LogoutHandler())
}
