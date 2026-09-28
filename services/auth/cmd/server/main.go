package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MaksimCpp/auth/internal/config"
	deliverygrpc "github.com/MaksimCpp/auth/internal/delivery/grpc"
	"github.com/MaksimCpp/auth/internal/delivery/grpc/middleware"
	jwtservice "github.com/MaksimCpp/auth/internal/infrastructure/jwt_service"
	repository "github.com/MaksimCpp/auth/internal/repository/postgresql"
	"github.com/MaksimCpp/auth/internal/usecase"
	authpb "github.com/MaksimCpp/auth/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	port := fmt.Sprintf(":%s", cfg.GRPCPort)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer pool.Close()

	tokenService := jwtservice.NewJWTTokenService(
		cfg.JWTSecretKey,
		15 * time.Minute,
		7 * 24 * time.Hour,
	)

	userRepo := repository.NewPostgreSQLUserRepository(pool)
	regiserUC := usecase.NewPostgreSQLRegisterUseCase(userRepo)
	loginUC := usecase.NewPostgreSQLLoginUseCase(userRepo, tokenService)

	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor(tokenService)))
	handler := deliverygrpc.NewAuthHandler(regiserUC, loginUC)
	authpb.RegisterAuthServiceServer(server, handler)

	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err.Error())
	}
}
