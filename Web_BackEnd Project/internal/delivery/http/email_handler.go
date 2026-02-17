package http

import (
	"mime/multipart"
	"net/http"

	"web_backend_project/internal/service"
)

// EmailHandler обрабатывает HTTP запросы связанные с email
type EmailHandler struct {
	emailService service.EmailServiceInterface
}

// NewEmailHandler создает новый экземпляр EmailHandler
func NewEmailHandler(emailService service.EmailServiceInterface) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// SendEmail обрабатывает запрос на отправку email
func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Ограничение загрузки 10 МБ
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	to := r.FormValue("to")
	subject := r.FormValue("subject")
	body := r.FormValue("body")

	// Обработка вложения файла
	var attachment multipart.File
	var filename string

	file, header, err := r.FormFile("attachment")
	if err == nil {
		// Файл найден
		attachment = file
		filename = header.Filename
		defer file.Close()
	} else if err != http.ErrMissingFile {
		// Ошибка при получении файла
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Отправляем email через сервис
	err = h.emailService.SendEmail(to, subject, body, attachment, filename)
	if err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
