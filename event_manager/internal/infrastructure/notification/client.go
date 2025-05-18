package notification

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gen "github.com/XRS0/ToTalkB/event_manager/internal/domain/gen"
)

type Client struct {
	client gen.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}

	client := gen.NewNotificationServiceClient(conn)
	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) SendNotification(ctx context.Context, notificationType string, payload []byte) (*gen.SendNotificationResponse, error) {
	resp, err := c.client.SendNotification(ctx, &gen.SendNotificationRequest{
		Type:    notificationType,
		Payload: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}
	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
