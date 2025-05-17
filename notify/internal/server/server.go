package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"notify/internal/config"
	grpcserver "notify/internal/infrastructure/grpc"
	"notify/internal/service"

	"google.golang.org/grpc"
)

type Server struct {
	config  *config.Config
	http    *http.Server
	grpc    *grpc.Server
	service *service.NotificationService
}

func NewServer(cfg *config.Config, svc *service.NotificationService) *Server {
	grpcServer := grpc.NewServer()
	grpcserver.RegisterServer(grpcServer, svc)

	srv := &Server{
		config:  cfg,
		service: svc,
		grpc:    grpcServer,
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", cfg.Server.Port),
		},
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.GRPCPort))
		if err != nil {
			panic(err)
		}
		if err := s.grpc.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.grpc.GracefulStop()
	return s.http.Shutdown(ctx)
}
