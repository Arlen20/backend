package http

import (
	"net/http"

	"github.com/rs/cors"
)

// NewRouter создает и настраивает все HTTP маршруты
func NewRouter(userHandler *UserHandler, emailHandler *EmailHandler) http.Handler {
	// Инициализация маршрутизатора
	mux := http.NewServeMux()

	// Маршруты для пользователей
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetUserByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.CreateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			userHandler.UpdateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			userHandler.DeleteUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Маршруты для email
	mux.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			emailHandler.SendEmail(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	return corsHandler
}
