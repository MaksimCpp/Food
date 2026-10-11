package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ProductsConfig struct {
	DBUrl        string

	GRPCPort     string
	GRPCAddress  string
}

func Load() *ProductsConfig {
	_ = godotenv.Load()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("PRODUCTS_DB")

	return &ProductsConfig{
		DBUrl: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password,
			host, port, db,
		),

		GRPCPort: os.Getenv("GRPC_PRODUCTS_PORT"),
		GRPCAddress: os.Getenv("GRPC_PRODUCTS_ADDRESS"),
	}
	
}