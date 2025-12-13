package auth

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// --- структура данных шаблона ---
type emailData struct {
	Email  string
	Status string
	Valid  bool
}

// --- функция валидации email ---
func validateEmail(email string) (bool, string) {
	if email == "" {
		log.Println("[validateEmail] ❌ Пустой email")
		return false, "⚠️ Email не указан"
	}

	// базовая проверка на email
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegexp.MatchString(email) {
		log.Println("[validateEmail] ❌ Недопустимый email")
		return false, "⚠️ Неверный формат email"
	}

	log.Println("[validateEmail] ✅ Email валиден")
	return true, "✅ Email принят"
}

// --- 1️⃣ Показываем форму ---
func EmailHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ShowEmailFormHandler] GET /email")

		data := emailData{
			Status: "Введите email@gmail.com",
			Valid:  false,
		}

		err := tmpl.ExecuteTemplate(w, "email", data)
		if err != nil {
			log.Printf("[ShowEmailFormHandler] ❌ Ошибка шаблона: %v", err)
		} else {
			log.Println("[ShowEmailFormHandler] ✅ Форма email загружена")
		}
	}
}

// --- 2️⃣ Проверка email ---
func EmailCheckHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ValidateEmailHandler] POST /validate-email")

		email := r.FormValue("email") // оставил "key" для совместимости с формой
		log.Printf("[ValidateEmailHandler] Получен email: %q\n", email)

		valid, msg := validateEmail(email)
		data := emailData{Email: email, Status: msg, Valid: valid}

		err := tmpl.ExecuteTemplate(w, "check__email", data)
		if err != nil {
			log.Printf("[ValidateEmailHandler] ❌ Ошибка шаблона: %v", err)
		} else {
			log.Println("[ValidateEmailHandler] ✅ email_check успешно отрендерен")
		}
	}
}
