package queue

import (
	"context"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	pb "github.com/k8sdeploy/protobufs/generated/queue_service/v1"
	ConfigBuilder "github.com/keloran/go-config"
)

type Server struct {
	pb.UnimplementedQueueServiceServer
	ConfigBuilder.Config
}

func (s *Server) CreateQueueAccount(ctx context.Context, in *pb.CreateQueueAccountRequest) (*pb.CreateQueueAccountResponse, error) {
	return &pb.CreateQueueAccountResponse{
		Status: pointerutil.StringPtr("Not Implemented"),
	}, nil
}

func (s *Server) DeleteQueueAccount(ctx context.Context, in *pb.DeleteQueueAccountRequest) (*pb.DeleteQueueAccountResponse, error) {
	return &pb.DeleteQueueAccountResponse{
		Status: pointerutil.StringPtr("Not Implemented"),
	}, nil
}

func (s *Server) GetQueueAccount(ctx context.Context, in *pb.GetQueueAccountRequest) (*pb.GetQueueAccountResponse, error) {
	return &pb.GetQueueAccountResponse{
		Status: pointerutil.StringPtr("Not Implemented"),
	}, nil
}
