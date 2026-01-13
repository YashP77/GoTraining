package main

import (
	"log"
	"log/slog"
	"net"

	"goTraining/api"
	pb "goTraining/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	filePath := "output/messages.txt"
	addr := ":50051"

	slog.Info("start gRPC server", "running on", addr, "storage", filePath)

	log.Printf("starting gRPC server on %s (storage: %s)\n", addr, filePath)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("listen failed", "err", err)
	}
	srv := grpc.NewServer()
	pb.RegisterMessageServiceServer(srv, api.NewServer(filePath))

	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		slog.Error("gRPC server failed", "err", err)
	}
}
