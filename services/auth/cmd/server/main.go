package main

import (
	"log"
	"net"

	deliverygrpc "github.com/MaksimCpp/auth/internal/delivery/grpc"
	authpb "github.com/MaksimCpp/auth/proto"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	server := grpc.NewServer()
	handler := deliverygrpc.NewAuthHandler()
	authpb.RegisterAuthServiceServer(server, handler)

	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err.Error())
	}
}
