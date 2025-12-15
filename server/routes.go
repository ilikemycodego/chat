package server

import (
	"chat/auth"
	"chat/middleware"
	"chat/setting"
	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует все маршруты
func RegisterRoutes(m *mux.Router, tmpl *template.Template) {

	// Основной маршрут

	m.Handle("/", middleware.UserContextMiddleware(auth.BaseHandler(tmpl)))
	m.HandleFunc("/start", auth.EmailHandler(tmpl))
	m.HandleFunc("/check-email", auth.EmailCheckHandler(tmpl))

	m.HandleFunc("/theme", auth.ToggleThemeHandler)

	m.HandleFunc("/get-password", auth.GetCodeHandler(tmpl))
	m.HandleFunc("/verify-code", auth.VerifyCodeHandler())

	// Пример дополнительных маршрутов:
	// r.HandleFunc("/login", auth.LoginHandler(tmpl))
	// r.HandleFunc("/logout", auth.LogoutHandler())

	m.HandleFunc("/setting", setting.SettingHandler(tmpl))
	m.HandleFunc("/user-setting", setting.UserSettingHandler(tmpl))
	m.HandleFunc("/name-setting", setting.NameHandler(tmpl))
	m.HandleFunc("/check-name", setting.NameCheckHandler(tmpl))

}
