package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
	DBUrl string

	GRPCPort string
}

func Load() *AuthConfig {
	_ = godotenv.Load()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("AUTH_DB")

	return &AuthConfig{
		DBUrl: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password,
			host, port, db,
		),

		GRPCPort: os.Getenv("GRPC_AUTH_PORT"),
	}
	
}