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
	m.HandleFunc("/start", auth.EmailHandler(tmpl))
	m.HandleFunc("/check-email", auth.EmailCheckHandler(tmpl))

	m.HandleFunc("/theme", auth.ToggleThemeHandler)

	m.HandleFunc("/get-password", auth.GetCodHandler(tmpl))
	m.HandleFunc("/verify-code", auth.VerifyCodeHandler(tmpl))

	// Пример дополнительных маршрутов:
	// r.HandleFunc("/login", auth.LoginHandler(tmpl))
	// r.HandleFunc("/logout", auth.LogoutHandler())
}
