package api

import (
	"context"
	"goTraining/internal"
	pb "goTraining/proto"
)

type server struct {
	pb.UnimplementedMessageServiceServer
	filePath string
}

func NewServer(filePath string) *server {
	return &server{filePath: filePath}
}

func (s *server) SaveMessage(ctx context.Context, req *pb.SaveMessageRequest) (*pb.SaveMessageResponse, error) {

	file := internal.OpenFile(ctx, s.filePath)
	defer file.Close()

	internal.WriteToFile(ctx, file, req.Message, int(req.UserID))
	return &pb.SaveMessageResponse{}, nil
}

func (s *server) GetLast10(ctx context.Context, _ *pb.GetLast10Request) (*pb.GetLast10Response, error) {

	lines := internal.ReadLastTen(ctx, s.filePath)
	return &pb.GetLast10Response{
		Messages: lines,
	}, nil
}
