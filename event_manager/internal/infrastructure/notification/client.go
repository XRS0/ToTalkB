package notification

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gen "github.com/XRS0/ToTalkB/event_manager/internal/domain/proto"

	"github.com/XRS0/ToTalkB/event_manager/internal/config"
)

type Client struct {
	client gen.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewClient(cfg *config.NotificationServiceConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %v", err)
	}

	client := gen.NewNotificationServiceClient(conn)
	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) SendNotification(ctx context.Context, notificationType string, payload []byte) (*gen.SendNotificationResponse, error) {
	req := &gen.SendNotificationRequest{
		Type:    notificationType,
		Payload: payload,
	}
	return c.client.SendNotification(ctx, req)
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
