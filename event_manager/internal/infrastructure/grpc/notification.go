package grpc

import (
	"context"

	"github.com/XRS0/ToTalkB/proto/gen"
	"google.golang.org/grpc"
)

type NotificationClient struct {
	client gen.NotificationServiceClient
}

func NewNotificationClient(conn *grpc.ClientConn) *NotificationClient {
	return &NotificationClient{
		client: gen.NewNotificationServiceClient(conn),
	}
}

func (c *NotificationClient) SendNotification(ctx context.Context, notificationType string, payload []byte) (string, error) {
	resp, err := c.client.SendNotification(ctx, &gen.SendNotificationRequest{
		Type:    notificationType,
		Payload: payload,
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *NotificationClient) GetNotificationStatus(ctx context.Context, id string) (string, error) {
	resp, err := c.client.GetNotificationStatus(ctx, &gen.GetNotificationStatusRequest{
		Id: id,
	})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
