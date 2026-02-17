package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"web_backend_project/internal/domain"
)

type mongoUserRepository struct {
	db         *mongo.Client
	database   string
	collection string
}

// NewMongoUserRepository creates a new instance of mongoUserRepository
func NewMongoUserRepository(db *mongo.Client, database, collection string) domain.UserRepository {
	return &mongoUserRepository{
		db:         db,
		database:   database,
		collection: collection,
	}
}

func (r *mongoUserRepository) GetUsers(ctx context.Context, page, limit int, filter, sortBy, sortOrder string) ([]domain.User, error) {
	collection := r.db.Database(r.database).Collection(r.collection)
	var users []domain.User

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

	cursor, err := collection.Find(ctx, filterQuery, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *mongoUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	collection := r.db.Database(r.database).Collection(r.collection)
	var user domain.User

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepository) CreateUser(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {
	collection := r.db.Database(r.database).Collection(r.collection)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.NilObjectID, fmt.Errorf("failed to get inserted ID")
}

func (r *mongoUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	collection := r.db.Database(r.database).Collection(r.collection)

	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoUserRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Database(r.database).Collection(r.collection)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
