package server

import (
	"context"
	"fmt"
	"net/http"

	"notify/internal/config"
	"notify/internal/infrastructure/kafka"
	"notify/internal/service"
)

type Server struct {
	config   *config.Config
	http     *http.Server
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

	srv := &Server{
		config:   cfg,
		service:  svc,
		consumer: consumer,
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", cfg.Server.Port),
		},
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	go s.consumer.Start(ctx)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
