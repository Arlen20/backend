package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

var client *mongo.Client

func main() {
	var err error
	mongoURI := "mongodb+srv://nurlybaynurbol:987412365nn@cluster0.436nq.mongodb.net/?retryWrites=true&w=majority"

	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	fmt.Println("MongoDB connection established")

	http.HandleFunc("/users", getUsers)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend's origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(http.DefaultServeMux)
	http.HandleFunc("/send-email", sendEmailHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
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

	collection := client.Database("test").Collection("users")
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
