package service

import (
  "fmt"
  "github.com/bugfixes/go-bugfixes/logs"
  pb "github.com/k8sdeploy/protobufs/generated/queue_service/v1"
  "github.com/k8sdeploy/queue-manager-service/internal/queue"
  ConfigBuilder "github.com/keloran/go-config"
  "github.com/keloran/go-healthcheck"
  "github.com/rs/cors"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
  "net"
  "net/http"
)

type Service struct {
	ConfigBuilder.Config
}

func NewService(cfg ConfigBuilder.Config) *Service {
	return &Service{
		Config: cfg,
	}
}

func (s *Service) Start() error {
	errChan := make(chan error)
	go startHTTP(s.Config, errChan)
	go startGRPC(s.Config, errChan)

	return <-errChan
}

func startHTTP(cfg ConfigBuilder.Config, errChan chan error) {
	logs.Local().Infof("Starting HTTP on %d", cfg.Local.HTTPPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthcheck.HTTP)
	handler := cors.Default().Handler(mux)
	errChan <- http.ListenAndServe(fmt.Sprintf(":%d", cfg.Local.HTTPPort), handler)
}

func startGRPC(cfg ConfigBuilder.Config, errChan chan error) {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Local.GRPCPort))
	if err != nil {
		errChan <- err
		return
	}

	gs := grpc.NewServer()
	reflection.Register(gs)
	pb.RegisterQueueServiceServer(gs, &queue.Server{
		Config: cfg,
	})
	logs.Local().Infof("Starting gRPC on %d", cfg.Local.GRPCPort)
	errChan <- gs.Serve(list)
}
