package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MaksimCpp/auth/internal/config"
	deliverygrpc "github.com/MaksimCpp/auth/internal/delivery/grpc"
	jwtservice "github.com/MaksimCpp/auth/internal/infrastructure/jwt_service"
	repository "github.com/MaksimCpp/auth/internal/repository/postgresql"
	"github.com/MaksimCpp/auth/internal/service"
	authpb "github.com/MaksimCpp/auth/proto"
	"github.com/MaksimCpp/pkg/interceptor"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	port := fmt.Sprintf(":%s", cfg.GRPCPort)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err.Error())
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	// defer cancel()

	pool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer pool.Close()

	tokenService := jwtservice.NewJWTTokenService(
		cfg.JWTSecretKey,
		15 * time.Minute,
		7 * 24 * time.Hour,
	)

	authConn, err := grpc.NewClient(
		cfg.GRPCAddress, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	authClient := authpb.NewAuthServiceClient(authConn)

	userRepo := repository.NewPostgreSQLUserRepository(pool)
	userService := service.NewPostgreSQLUserService(cfg.AdminCode, userRepo, tokenService)

	var publicMethods = map[string]struct{}{
		"/auth.AuthService/RegisterUser":        {},
		"/auth.AuthService/RegisterAdmin":       {},
		"/auth.AuthService/LoginUser":           {},
		"/auth.AuthService/LoginAdmin":          {},
		"/auth.AuthService/LoginCourier":        {},
		"/auth.AuthService/Refresh":             {},
		"/auth.AuthService/ValidateAccessToken": {},
	}

	var methodsForAdmin = map[string]struct{} {
		"/auth.AuthService/RegisterCourier": {},
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.AuthInterceptor(authClient, publicMethods),
			interceptor.RoleInterceptor(methodsForAdmin),
		),
	)
	handler := deliverygrpc.NewAuthHandler(userService, tokenService)
	authpb.RegisterAuthServiceServer(server, handler)

	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err.Error())
	}
}
