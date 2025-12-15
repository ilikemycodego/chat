package setting

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// --- структура данных шаблона ---
type nameData struct {
	Name   string
	Status string
	Valid  bool
}

// --- функция валидации name ---
func validateName(name string) (bool, string) {
	if name == "" {
		log.Println("[validateName] ❌ Пустое имя")
		return false, "⚠️ Имя не указано"
	}

	// пример базовой проверки (можно убрать проверку на email-формат)
	nameRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !nameRegexp.MatchString(name) {
		log.Println("[validateName] ❌ Недопустимое имя")
		return false, "⚠️ Неверный формат имени"
	}

	log.Println("[validateName] ✅ Имя норм")
	return true, "✅ Name принят"
}

// --- 1️⃣ Показываем форму ---
func NameHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ShowNameFormHandler] GET /name")

		data := nameData{
			Status: "Введите ваше имя",
			Valid:  false,
		}

		err := tmpl.ExecuteTemplate(w, "name__setting", data)
		if err != nil {
			log.Printf("[ShowNameFormHandler] ❌ Ошибка шаблона: %v", err)
		} else {
			log.Println("[ShowNameFormHandler] ✅ Форма имени загружена")
		}
	}
}

// --- 2️⃣ Проверка name ---
func NameCheckHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ValidateNameHandler] POST /validate-name")

		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name") // оставлено для совместимости с формой
		log.Printf("[ValidateNameHandler] Получено имя: %q\n", name)

		valid, msg := validateName(name)
		data := nameData{Name: name, Status: msg, Valid: valid}

		err := tmpl.ExecuteTemplate(w, "check__name", data)
		if err != nil {
			log.Printf("[ValidateNameHandler] ❌ Ошибка шаблона: %v", err)
		} else {
			log.Println("[ValidateNameHandler] ✅ name_check успешно отрендерен")
		}
	}
}
