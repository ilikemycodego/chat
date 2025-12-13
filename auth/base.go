package auth

import (
	"html/template"
	"log"
	"net/http"
)

// BaseHandler рендерит шаблон с логированием ошибок
func BaseHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		theme := "light"
		if c, err := r.Cookie("theme"); err == nil && c.Value == "dark" {
			theme = "dark"
		}

		data := struct{ Theme string }{Theme: theme}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Printf("Ошибка рендеринга шаблона base.html для маршрута %s: %v", r.URL.Path, err)
			http.Error(w, "Ошибка сервера при загрузке шаблона", http.StatusInternalServerError)
		}
	}
}

// переключение темы через cookie
func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// читаем текущую тему
	theme := "light"
	if c, err := r.Cookie("theme"); err == nil && c.Value == "dark" {
		theme = "dark"
	}

	// переключаем тему
	if theme == "light" {
		theme = "dark"
	} else {
		theme = "light"
	}

	// ставим cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    theme,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   30 * 24, // 1 месяц
	})

	// редирект для HTMX
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
