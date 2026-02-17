package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	HTTPServerPort string
	GRPCServerPort string

	// MongoDB configuration
	MongoURI        string
	MongoDB         string
	MongoCollection string

	// Redis configuration
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	CacheTTL      int // в секундах

	// Other configurations
	Debug bool
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Parse Redis DB number
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	// Parse Cache TTL
	cacheTTL, _ := strconv.Atoi(getEnv("CACHE_TTL", "300")) // 5 минут по умолчанию

	config := &Config{
		// Server configuration with defaults
		HTTPServerPort: getEnv("HTTP_SERVER_PORT", "8080"),
		GRPCServerPort: getEnv("GRPC_SERVER_PORT", "50051"),

		// MongoDB configuration with defaults
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "test"),
		MongoCollection: getEnv("MONGO_COLLECTION", "users"),

		// Redis configuration with defaults
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
		CacheTTL:      cacheTTL,

		// Other configurations with defaults
		Debug: getEnvAsBool("DEBUG", false),
	}

	return config
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Invalid boolean value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}

	return boolValue
}
