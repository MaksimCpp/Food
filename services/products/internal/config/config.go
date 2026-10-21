package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ProductsConfig struct {
	DBUrl            string

	GRPCPort         string
	GRPCAddress      string
	GRPCAuthAddress  string

	// MinioApiPort     string
	// MinioConsolePort string
	MinioEndpoint    string
	MinioBucket      string
	MinioUseSSL      bool
	MinioAccessKey   string
	MinioSecretKey   string

	RedisAddress     string
	RedisTTLMinutes  int
}

func Load() *ProductsConfig {
	_ = godotenv.Load()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("PRODUCTS_DB")

	redisTTL, err := strconv.Atoi(os.Getenv("REDIS_TTL_MINUTES"))

	if err != nil {
		return nil
	}

	return &ProductsConfig{
		DBUrl: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password,
			host, port, db,
		),

		GRPCPort: os.Getenv("GRPC_PRODUCTS_PORT"),
		GRPCAddress: os.Getenv("GRPC_PRODUCTS_ADDRESS"),
		GRPCAuthAddress: os.Getenv("GRPC_AUTH_ADDRESS"),

		// MinioApiPort: os.Getenv("MINIO_API_PORT"),
		// MinioConsolePort: os.Getenv("MINIO_CONSOLE_PORT"),
		MinioEndpoint: os.Getenv("MINIO_ENDPOINT"),
		MinioBucket: os.Getenv("MINIO_BUCKET"),
		MinioUseSSL: os.Getenv("MINIO_USE_SSL") == "true",
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),

		RedisAddress: os.Getenv("REDIS_ADDRESS"),
		RedisTTLMinutes: redisTTL,
	}
	
}