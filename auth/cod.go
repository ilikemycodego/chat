package auth

import (
	"chat/db"
	"context"
	"fmt"
	"html/template"
	"log"
	"math/big"

	"net/http"
	"os"
	"time"

	"crypto/rand"

	"github.com/mailersend/mailersend-go"
)

type EmailData struct {
	Email  string
	Status string
	Valid  bool
}

func generateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

func canSendCode(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := db.DB.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM auth_codes
		WHERE email=$1 AND created_at > NOW() - INTERVAL '1 hour'
	`, email).Scan(&count)

	if err != nil {
		return false, err
	}

	return count < 10, nil // максимум 5 кодов в час
}

// Сохраняем код для email в таблицу auth_codes
func saveCode(email, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.DB.Exec(ctx, `
		INSERT INTO auth_codes (email, code, created_at, used)
		VALUES ($1, $2, $3, $4)
	`, email, code, time.Now(), false)

	return err
}

// Отправка кода на email
func sendCode(email, code string) error {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")

	if apiKey == "" || fromEmail == "" {
		return fmt.Errorf("MAILERSEND_API_KEY или FROM_EMAIL не заданы")
	}

	log.Println("[sendCode] MailerSend инициализирован")
	log.Println("[sendCode] Отправитель:", fromEmail)
	log.Println("[sendCode] Получатель:", email)

	ms := mailersend.NewMailersend(apiKey)
	ctx := context.Background()

	from := mailersend.From{
		Email: fromEmail,
		Name:  "Your Service",
	}

	recipients := []mailersend.Recipient{
		{Email: email},
	}

	message := ms.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject("Ваш регистрационный код")
	message.SetText("Ваш код для входа: " + code)

	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		log.Println("[sendCode] ❌ Ошибка отправки письма:", err)
		return err
	}

	log.Println("[sendCode] ✅ Письмо отправлено, X-Message-Id:", res.Header.Get("X-Message-Id"))
	return nil
}

// ------------------------------------------------

//-------------------------

func GetCodeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Получаем email из формы
		email := r.FormValue("email")
		if email == "" {
			fmt.Fprint(w, "Введите email!")
			return
		}

		// Проверяем, можно ли отправлять код
		ok, err := canSendCode(email)
		if err != nil {
			http.Error(w, "Ошибка сервера", 500)
			return
		}
		if !ok {
			fmt.Fprint(w, "Слишком много запросов, попробуйте позже")
			return
		}

		// Генерация кода
		code := generateCode()

		// Сохраняем в БД
		if err := saveCode(email, code); err != nil {
			log.Println("[GetPasswordHandler] ❌ Ошибка сохранения кода:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Отправка письма
		if err := sendCode(email, code); err != nil {
			log.Println("[GetPasswordHandler] ❌ Ошибка отправки письма:", err)
			fmt.Fprint(w, "Не удалось отправить код, попробуйте позже")
			return
		}

		// --- рендерим форму для ввода кода/пароля ---
		if err := tmpl.ExecuteTemplate(w, "get__code", map[string]string{"Email": email}); err != nil {
			log.Printf("[GetPasswordHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}

		log.Println("[GetPasswordHandler] ✅ код отправлен и форма для ввода пароля отрендерена")
	}
}

//---------------
