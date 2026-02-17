package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"web_backend_project/internal/domain"
	"web_backend_project/pkg/cache"
)

type userUseCase struct {
	userRepo    domain.UserRepository
	redisClient *cache.RedisClient
	cacheTTL    time.Duration
}

// NewUserUseCase creates a new instance of userUseCase
func NewUserUseCase(userRepo domain.UserRepository, redisClient *cache.RedisClient, cacheTTL int) domain.UserUseCase {
	return &userUseCase{
		userRepo:    userRepo,
		redisClient: redisClient,
		cacheTTL:    time.Duration(cacheTTL) * time.Second,
	}
}

func (u *userUseCase) GetUsers(ctx context.Context, page, limit int, filter, sortBy, sortOrder string) ([]domain.User, error) {
	return u.userRepo.GetUsers(ctx, page, limit, filter, sortBy, sortOrder)
}

func (u *userUseCase) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("user:%s", id.Hex())

	// Проверяем наличие данных в кэше
	var cachedUser domain.User
	err := u.redisClient.Get(ctx, cacheKey, &cachedUser)
	if err == nil {
		log.Printf("Cache hit for key: %s", cacheKey)
		return &cachedUser, nil
	}
	log.Printf("Cache miss for key: %s, error: %v", cacheKey, err)

	// Если данных нет в кэше, получаем из репозитория
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.Printf("User with ID %s not found in database", id.Hex())
		return nil, nil
	}

	// Сохраняем данные в кэш
	if err := u.redisClient.Set(ctx, cacheKey, user, u.cacheTTL); err != nil {
		log.Printf("Failed to cache user with key %s: %v", cacheKey, err)
	} else {
		log.Printf("User cached successfully with key: %s", cacheKey)
	}

	return user, nil
}

func (u *userUseCase) CreateUser(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {
	// Здесь можно добавить бизнес-логику/валидацию перед созданием пользователя
	return u.userRepo.CreateUser(ctx, user)
}

func (u *userUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
	// Обновляем пользователя
	err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	// Инвалидируем кэш
	cacheKey := fmt.Sprintf("user:%s", user.ID.Hex())
	if err := u.redisClient.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate user cache for key %s: %v", cacheKey, err)
	} else {
		log.Printf("Cache invalidated for key: %s", cacheKey)
	}

	return nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	// Удаляем пользователя
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// Инвалидируем кэш
	cacheKey := fmt.Sprintf("user:%s", id.Hex())
	if err := u.redisClient.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate user cache for key %s: %v", cacheKey, err)
	} else {
		log.Printf("Cache invalidated for key: %s", cacheKey)
	}

	return nil
}
