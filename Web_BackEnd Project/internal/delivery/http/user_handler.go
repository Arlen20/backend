package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"web_backend_project/internal/domain"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetUsers handles GET /users request
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := 1
	limit := 5
	filter := ""
	sortBy := "firstName"
	sortOrder := "asc"

	queries := r.URL.Query()
	if pageParam := queries.Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}
	if limitParam := queries.Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
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

	// Get users from use case
	users, err := h.userUseCase.GetUsers(context.Background(), page, limit, filter, sortBy, sortOrder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching users: %v", err), http.StatusInternalServerError)
		return
	}

	// Return users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID handles GET /users/:id request
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL or request parameters as needed
	// This will depend on your router implementation
	idStr := r.URL.Query().Get("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetUserByID(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser handles POST /users request
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert map to domain.User
	var user domain.User
	userBytes, _ := json.Marshal(userData)
	json.Unmarshal(userBytes, &user)

	id, err := h.userUseCase.CreateUser(context.Background(), &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id.Hex(),
		"message": "User created successfully",
	})
}

// UpdateUser handles PUT /users request
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData bson.M
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if ID exists in request
	idInterface, exists := userData["_id"]
	if !exists {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID to ObjectID
	idStr, ok := idInterface.(string)
	if !ok {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Convert map to domain.User
	var user domain.User
	userBytes, _ := json.Marshal(userData)
	json.Unmarshal(userBytes, &user)
	user.ID = id

	if err := h.userUseCase.UpdateUser(context.Background(), &user); err != nil {
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

// DeleteUser handles DELETE /users request
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.DeleteUser(context.Background(), id); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}
