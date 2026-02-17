package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"

	"web_backend_project/grpc"
	"web_backend_project/pkg/cache"
	"web_backend_project/quiz"
	"web_backend_project/transaction"
)

var mainClient *mongo.Client
var redisClient *cache.RedisClient
var cacheTTL time.Duration
var nc *nats.Conn // NATS connection

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Загрузка настроек Redis
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0           // По умолчанию используем DB 0
	cacheTTLSeconds := 300 // По умолчанию 5 минут

	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	mainClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения к MongoDB
	err = mainClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection established")

	// Connect to NATS server
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS server")

	// Subscribe to NATS subject for email notifications
	go subscribeToEmailNotifications()

	// Инициализация Redis
	redisClient, err = cache.NewRedisClient(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without caching.", err)
	} else {
		fmt.Println("Redis connection established")
		defer redisClient.Close()
	}

	// Установка TTL для кэша
	cacheTTL = time.Duration(cacheTTLSeconds) * time.Second

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	// Настройка маршрутов
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/users/get", getUserByID) // Добавляем новый обработчик с поддержкой кэширования
	http.HandleFunc("/users/create", createUser)
	http.HandleFunc("/users/update", updateUser)
	http.HandleFunc("/users/delete", deleteUser)
	http.HandleFunc("/send-email", sendEmailHandler)
	http.HandleFunc("/demo/cache", cacheDemoHandler) // Эндпоинт для демонстрации кэширования

	// Запускаем все сервисы
	go func() {
		if err := grpc.StartGRPCServer(50051); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
	go transaction.StartTransactionService()
	go quiz.StartQuizService()

	// Запускаем основной сервер
	log.Printf("Starting main server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Failed to start main server: %v", err)
	}
}

// getUserByID получает пользователя по ID с использованием кэширования Redis
func getUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Преобразуем строку ID в ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("user:%s", idStr)

	var user map[string]interface{}

	// Проверяем наличие данных в кэше, если Redis доступен
	if redisClient != nil {
		// Пытаемся получить данные из кэша
		err := redisClient.Get(r.Context(), cacheKey, &user)
		if err == nil {
			// Данные найдены в кэше, возвращаем их
			fmt.Println("Cache hit for user:", idStr)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(user)
			return
		}
		fmt.Println("Cache miss for user:", idStr, "Error:", err)
	}

	// Данных нет в кэше или Redis недоступен, получаем из MongoDB
	collection := mainClient.Database("test").Collection("users")
	err = collection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Сохраняем данные в кэш, если Redis доступен
	if redisClient != nil {
		if err := redisClient.Set(r.Context(), cacheKey, user, cacheTTL); err != nil {
			// Логируем ошибку, но продолжаем работу
			fmt.Printf("Failed to cache user: %v\n", err)
		} else {
			fmt.Println("User cached with key:", cacheKey)
		}
	}

	// Возвращаем данные клиенту
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(user)
}

// cacheDemoHandler демонстрирует работу кэша
func cacheDemoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID пользователя из параметров запроса
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("user:%s", userIDStr)

	// Проверяем, есть ли данные в кэше
	var cachedUser map[string]interface{}
	cacheHit := false

	if redisClient != nil {
		err := redisClient.Get(r.Context(), cacheKey, &cachedUser)
		cacheHit = err == nil
	}

	// Информация о состоянии кэша
	cacheInfo := map[string]interface{}{
		"cache_hit": cacheHit,
		"cache_key": cacheKey,
	}

	if cacheHit {
		// Если данные есть в кэше, удаляем их для демонстрации обоих сценариев
		if redisClient != nil {
			// Удаляем ключ из кэша
			if err := redisClient.Delete(r.Context(), cacheKey); err != nil {
				fmt.Printf("Failed to delete key from cache: %v\n", err)
			} else {
				fmt.Printf("Key %s deleted from cache for demo\n", cacheKey)
			}
		}
	}

	// Измеряем время первого запроса (без кэша)
	start := time.Now()

	// Получаем пользователя из БД
	collection := mainClient.Database("test").Collection("users")
	var user map[string]interface{}
	err = collection.FindOne(r.Context(), bson.M{"_id": userID}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Время запроса без кэша
	dbQueryTime := time.Since(start)

	// Сохраняем в кэш
	if redisClient != nil {
		if err := redisClient.Set(r.Context(), cacheKey, user, cacheTTL); err != nil {
			fmt.Printf("Failed to cache user for demo: %v\n", err)
		} else {
			fmt.Printf("User cached with key %s for demo\n", cacheKey)
		}
	}

	// Измеряем время второго запроса (с кэшем)
	start = time.Now()

	if redisClient != nil {
		var cachedUser map[string]interface{}
		err = redisClient.Get(r.Context(), cacheKey, &cachedUser)
		if err != nil {
			fmt.Printf("Cache miss after setting: %v\n", err)
		}
	}

	// Время запроса с кэшем
	cacheQueryTime := time.Since(start)

	// Формируем ответ с информацией о запросах
	response := map[string]interface{}{
		"user": user,
		"cache_demo": map[string]interface{}{
			"db_query": map[string]interface{}{
				"duration_ms": dbQueryTime.Milliseconds(),
				"cached":      false,
			},
			"cache_query": map[string]interface{}{
				"duration_ms": cacheQueryTime.Milliseconds(),
				"cached":      true,
			},
			"cache_info": cacheInfo,
			"speed_improvement": map[string]interface{}{
				"factor":         float64(dbQueryTime) / float64(cacheQueryTime),
				"percent_faster": (1 - float64(cacheQueryTime)/float64(dbQueryTime)) * 100,
			},
		},
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 5
	filter := ""
	sortBy := "name"
	sortOrder := "asc"

	queries := r.URL.Query()
	if pageParam := queries.Get("page"); pageParam != "" {
		fmt.Sscanf(pageParam, "%d", &page)
	}
	if limitParam := queries.Get("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}
	if filterParam := queries.Get("filter"); filterParam != "" {
		filter = filterParam
	}
	if sortByParam := queries.Get("sort_by"); sortByParam != "" {
		sortBy = sortByParam
	}
	if sortOrderParam := queries.Get("sort_order"); sortOrderParam != "" {
		sortOrder = sortOrderParam
	}

	fmt.Printf("Page: %d, Limit: %d, Filter: %s, Sort By: %s, Sort Order: %s\n", page, limit, filter, sortBy, sortOrder)

	collection := mainClient.Database("test").Collection("users")
	var users []map[string]interface{}

	options := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))

	if sortOrder == "desc" {
		options.SetSort(bson.D{{Key: sortBy, Value: -1}})
	} else {
		options.SetSort(bson.D{{Key: sortBy, Value: 1}})
	}

	filterQuery := bson.M{}
	if filter != "" {
		filterQuery = bson.M{
			"$or": []bson.M{
				{"firstName": bson.M{"$regex": filter, "$options": "i"}},
				{"lastName": bson.M{"$regex": filter, "$options": "i"}},
				{"username": bson.M{"$regex": filter, "$options": "i"}},
			},
		}
	}

	fmt.Printf("Executing MongoDB query with filter: %+v\n", filterQuery)
	cursor, err := collection.Find(context.Background(), filterQuery, options)
	if err != nil {
		log.Println("Error querying MongoDB:", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user map[string]interface{}
		if err := cursor.Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error iterating cursor:", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	to := r.FormValue("to")
	subject := r.FormValue("subject")
	body := r.FormValue("body")

	// Handle file attachment
	file, header, err := r.FormFile("attachment")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Configure email
	m := gomail.NewMessage()
	m.SetHeader("From", "nurlybaynurbol@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if file != nil {
		m.Attach(header.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, file)
			return err
		}))
	}

	// Send email
	d := gomail.NewDialer("smtp.gmail.com", 587, "nurlybaynurbol@gmail.com", "rdhk amua afhc mivw")

	if err := d.DialAndSend(m); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user bson.M
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := mainClient.Database("test").Collection("users")
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": result.InsertedID,
	})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user bson.M
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := user["_id"]
	delete(user, "_id")

	collection := mainClient.Database("test").Collection("users")
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": user},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Инвалидация кэша при обновлении
	if redisClient != nil && id != nil {
		cacheKey := fmt.Sprintf("user:%v", id)
		if err := redisClient.Delete(context.Background(), cacheKey); err != nil {
			// Логируем ошибку, но продолжаем работу
			fmt.Printf("Failed to invalidate user cache: %v\n", err)
		} else {
			fmt.Printf("Cache invalidated for key: %s\n", cacheKey)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"matchedCount":  result.MatchedCount,
		"modifiedCount": result.ModifiedCount,
	})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := mainClient.Database("test").Collection("users")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Инвалидация кэша при удалении
	if redisClient != nil {
		cacheKey := fmt.Sprintf("user:%s", idStr)
		if err := redisClient.Delete(context.Background(), cacheKey); err != nil {
			// Логируем ошибку, но продолжаем работу
			fmt.Printf("Failed to invalidate user cache: %v\n", err)
		} else {
			fmt.Printf("Cache invalidated for key: %s\n", cacheKey)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deletedCount": result.DeletedCount,
	})
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func sendEmail(to, subject, body string) {
	from := getEnv("EMAIL_USER", "")
	password := getEnv("EMAIL_PASS", "")
	if from == "" || password == "" {
		log.Println("EMAIL_USER or EMAIL_PASS is not set; skipping email send")
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending notification email:", err)
	}
}

func subscribeToEmailNotifications() {
	_, err := nc.Subscribe("email.notifications", func(msg *nats.Msg) {
		var emailData map[string]string
		if err := json.Unmarshal(msg.Data, &emailData); err != nil {
			log.Println("Error unmarshalling NATS message:", err)
			return
		}

		// Process the email data (e.g., send an email)
		log.Printf("Received email notification: %+v\n", emailData)
		sendEmail(emailData["to"], emailData["subject"], emailData["body"])
	})
	if err != nil {
		log.Println("Error subscribing to NATS subject:", err)
	}
}
