package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	authpb "github.com/MaksimCpp/auth/proto"
	"github.com/MaksimCpp/pkg/interceptor"
	"github.com/MaksimCpp/products/internal/config"
	deliverygrpc "github.com/MaksimCpp/products/internal/delivery/grpc"
	miniostorage "github.com/MaksimCpp/products/internal/infrastructure/minio"
	rediscache "github.com/MaksimCpp/products/internal/infrastructure/redis"
	repository "github.com/MaksimCpp/products/internal/repository/postgresql"
	"github.com/MaksimCpp/products/internal/service"
	productspb "github.com/MaksimCpp/products/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DBUrl)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer pool.Close()

	productsRepo := repository.NewPostgreSQLProductsRepository(pool)
	minioStorage, err := miniostorage.NewMinioStorage(
		cfg.MinioEndpoint,
		cfg.MinioPublicEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket, cfg.MinioUseSSL,
	)
	cache := rediscache.NewRedisCache(cfg.RedisAddress, time.Duration(cfg.RedisTTLMinutes))

	if err != nil {
		log.Fatal(err.Error())
	}

	productsService := service.NewPostgreSQLProductsService(
		productsRepo, minioStorage, cache,
	)

	handler := deliverygrpc.NewProductsHandler(productsService)

	authConn, err := grpc.NewClient(
		cfg.GRPCAuthAddress, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer authConn.Close()

	authClient := authpb.NewAuthServiceClient(authConn)

	publicMethods := map[string]struct{}{
		"/products.ProductsService/GetProduct":   {},
		"/products.ProductsService/ListProducts": {},
	}

	methodsForAdmin := map[string]struct{}{
		"/products.ProductsService/CreateProduct": {},
		"/products.ProductsService/UpdateProduct": {},
		"/products.ProductsService/DeleteProduct": {},
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.AuthInterceptor(authClient, publicMethods),
			interceptor.RoleInterceptor(methodsForAdmin),
		),
	)
	productspb.RegisterProductsServiceServer(server, handler)

	port := fmt.Sprintf(":%s", cfg.GRPCPort)

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = server.Serve(lis)

	if err != nil {
		log.Fatal(err.Error())
	}
}
