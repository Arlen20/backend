package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// UserRepository represents the user repository contract
type UserRepository interface {
	GetUsers(ctx context.Context, page, limit int, filter, sortBy, sortOrder string) ([]User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	CreateUser(ctx context.Context, user *User) (primitive.ObjectID, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}

// UserUseCase represents the user use case contract
type UserUseCase interface {
	GetUsers(ctx context.Context, page, limit int, filter, sortBy, sortOrder string) ([]User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	CreateUser(ctx context.Context, user *User) (primitive.ObjectID, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}
