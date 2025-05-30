package grpc

import (
	"context"
	"time"

	"github.com/XRS0/ToTalkB/notify/internal/domain/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client gen.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := gen.NewNotificationServiceClient(conn)
	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendNotification(ctx context.Context, userID int, notificationType string, payload []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.SendNotification(ctx, &gen.SendNotificationRequest{
		Type:    notificationType,
		Payload: payload,
	})
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}

func (c *Client) GetNotificationStatus(ctx context.Context, id string) (*gen.GetNotificationStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.GetNotificationStatus(ctx, &gen.GetNotificationStatusRequest{
		Id: id,
	})
}
