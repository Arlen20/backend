package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

var client *mongo.Client
var dbName = "test"
var transactionCollection = "transactions"
var usersCollection = "users"

type CartItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TransactionRequest struct {
	CartItems []CartItem `json:"cartItems"`
	Customer  Customer   `json:"customer"`
}

type PaymentForm struct {
	TransactionID  string `json:"transactionID"`
	CardNumber     string `json:"cardNumber"`
	ExpirationDate string `json:"expirationDate"`
	CVV            string `json:"cvv"`
	Name           string `json:"name"`
	Address        string `json:"address"`
}

type Transaction struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CartItems  []CartItem         `bson:"cartItems"`
	Customer   Customer           `bson:"customer"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	ReceiptURL string             `bson:"receipt_url"`
}

func startTransactionService() {
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

	http.HandleFunc("/transaction", handleTransaction)
	http.HandleFunc("/payment", handlePayment)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend's origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(http.DefaultServeMux)
	fmt.Println("Transaction microservice started at :8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var transactionRequest TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding transaction request:", err)
		return
	}

	transaction := Transaction{
		CartItems: transactionRequest.CartItems,
		Customer:  transactionRequest.Customer,
		Status:    "Pending Payment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := client.Database(dbName).Collection(transactionCollection)
	res, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		log.Println("Error inserting transaction into MongoDB:", err)
		return
	}

	transactionID := res.InsertedID.(primitive.ObjectID).Hex()
	paymentForm := PaymentForm{
		TransactionID:  transactionID,
		CardNumber:     "1234 5678 9012 3456",
		ExpirationDate: "12/24",
		CVV:            "123",
		Name:           transactionRequest.Customer.Name,
		Address:        "123 Main Street",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentForm)
}

func handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var paymentForm PaymentForm
	if err := json.NewDecoder(r.Body).Decode(&paymentForm); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding payment form:", err)
		return
	}

	// Mock payment processing
	paymentSuccess := paymentForm.TransactionID != "fail"

	// Update transaction status based on payment result
	status := "Declined"
	if paymentSuccess {
		status = "Paid"
	}

	transactionID, err := primitive.ObjectIDFromHex(paymentForm.TransactionID)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		log.Println("Error converting transaction ID:", err)
		return
	}

	collection := client.Database(dbName).Collection(transactionCollection)
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": transactionID}, update)
	if err != nil {
		http.Error(w, "Failed to update transaction status", http.StatusInternalServerError)
		log.Println("Error updating transaction status in MongoDB:", err)
		return
	}

	if paymentSuccess {
		// Generate and send fiscal receipt
		receiptURL := generateFiscalReceipt(transactionID.Hex(), paymentForm)
		update = bson.M{"$set": bson.M{"receipt_url": receiptURL, "status": "Completed"}}
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": transactionID}, update)
		if err != nil {
			http.Error(w, "Failed to update transaction with receipt URL", http.StatusInternalServerError)
			log.Println("Error updating transaction with receipt URL in MongoDB:", err)
			return
		}

		// Fetch transaction to get the customer ID
		var transaction Transaction
		err = collection.FindOne(context.Background(), bson.M{"_id": transactionID}).Decode(&transaction)
		if err != nil {
			log.Println("Error finding transaction:", err)
			return
		}

		// Fetch user email from MongoDB using the customer ID
		var customer Customer
		customerID, err := primitive.ObjectIDFromHex(transaction.Customer.ID)
		if err != nil {
			log.Println("Error converting customer ID:", err)
			return
		}
		err = client.Database(dbName).Collection(usersCollection).FindOne(context.Background(), bson.M{"_id": customerID}).Decode(&customer)
		if err != nil {
			log.Println("Error finding customer:", err)
			return
		}

		// Update the user's role to admin
		userCollection := client.Database(dbName).Collection(usersCollection)
		_, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": customerID}, bson.M{"$set": bson.M{"role": "admin"}})
		if err != nil {
			log.Println("Error updating user role:", err)
			return
		}

		sendReceiptEmail(customer.Email, receiptURL)
	}

	w.Header().Set("Content-Type", "application/json")
	if paymentSuccess {
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Transaction %s", status)})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to process payment"})
	}
}

func generateFiscalReceipt(transactionID string, paymentForm PaymentForm) string {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Company/Project name
	pdf.Cell(40, 10, "Company/Project Name: XYZ Inc.")
	pdf.Ln(12)

	// Transaction ID
	pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", transactionID))
	pdf.Ln(12)

	// Order date and time
	pdf.Cell(40, 10, fmt.Sprintf("Order Date: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(12)

	// Items
	pdf.Cell(40, 10, "Items:")
	pdf.Ln(12)

	collection := client.Database(dbName).Collection(transactionCollection)
	var transaction Transaction
	objID, _ := primitive.ObjectIDFromHex(transactionID)
	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&transaction)
	if err != nil {
		log.Println("Error finding transaction:", err)
		return ""
	}

	for _, item := range transaction.CartItems {
		pdf.Cell(40, 10, fmt.Sprintf("Product: %s - Price: $%.2f - Quantity: %d", item.Name, item.Price, item.Quantity))
		pdf.Ln(6)
	}

	// Total amount
	pdf.Cell(40, 10, fmt.Sprintf("Total Amount: $%.2f", calculateTotal(transaction.CartItems)))
	pdf.Ln(12)

	// Customer's full name
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", transaction.Customer.Name))
	pdf.Ln(12)

	// Payment method with encrypted card details
	pdf.Cell(40, 10, fmt.Sprintf("Payment Method: **** **** **** %s", paymentForm.CardNumber[len(paymentForm.CardNumber)-4:]))
	pdf.Ln(12)

	fileName := fmt.Sprintf("receipt_%s.pdf", transactionID)
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		log.Println("Error generating PDF:", err)
		return ""
	}
	return fileName
}

func calculateTotal(cartItems []CartItem) float64 {
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func sendReceiptEmail(to, receiptURL string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "nurlybaynurbol@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Fiscal Receipt")
	m.SetBody("text/plain", "Thank you for your purchase! Please find your receipt attached.")
	m.Attach(receiptURL)

	d := gomail.NewDialer("smtp.gmail.com", 587, "nurlybaynurbol@gmail.com", "rdhk amua afhc mivw")
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending receipt email:", err)
	}
}
