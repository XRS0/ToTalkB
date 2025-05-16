package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"notify/internal/config"
	grpcserver "notify/internal/infrastructure/grpc"
	"notify/internal/infrastructure/kafka"
	"notify/internal/service"

	"google.golang.org/grpc"
)

type Server struct {
	config   *config.Config
	http     *http.Server
	grpc     *grpc.Server
	service  *service.NotificationService
	consumer *kafka.Consumer
}

func NewServer(cfg *config.Config) *Server {
	kafkaConfig := kafka.ConsumerConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topics: kafka.Topics{
			Notifications: cfg.Kafka.Topics.Notifications,
		},
	}

	svc := service.NewNotificationService()
	consumer := kafka.NewConsumer(kafkaConfig, svc)
	grpcServer := grpc.NewServer()
	grpcserver.RegisterServer(grpcServer, svc)

	srv := &Server{
		config:   cfg,
		service:  svc,
		consumer: consumer,
		grpc:     grpcServer,
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

	go s.consumer.Start(ctx)

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.grpc.GracefulStop()
	return s.http.Shutdown(ctx)
}
