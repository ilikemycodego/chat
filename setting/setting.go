package setting

import (
	"html/template"
	"log"
	"net/http"
)

// SettingHandler рендерит шаблон с настройками
func SettingHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ОДИН рендер
		if err := tmpl.ExecuteTemplate(w, "setting", nil); err != nil {
			log.Printf("[SettingHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}

		log.Println("[SettingHandler] ✅ Страница отрендерена")
	}
}

func UserSettingHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ОДИН рендер
		if err := tmpl.ExecuteTemplate(w, "user__setting", nil); err != nil {
			log.Printf("[UserSettingHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}

		log.Println("[UserSettingHandler] ✅ Страница отрендерена")
	}
}
