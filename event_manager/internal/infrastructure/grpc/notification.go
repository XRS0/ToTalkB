package grpc

import (
	"context"

	"github.com/XRS0/ToTalkB/proto/gen_notify"
	"google.golang.org/grpc"
)

type NotificationClient struct {
	client gen_notify.NotificationServiceClient
}

func NewNotificationClient(conn *grpc.ClientConn) *NotificationClient {
	return &NotificationClient{
		client: gen_notify.NewNotificationServiceClient(conn),
	}
}

func (c *NotificationClient) SendNotification(ctx context.Context, notificationType string, payload []byte) (string, error) {
	resp, err := c.client.SendNotification(ctx, &gen_notify.SendNotificationRequest{
		Type:    notificationType,
		Payload: payload,
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *NotificationClient) GetNotificationStatus(ctx context.Context, id string) (string, error) {
	resp, err := c.client.GetNotificationStatus(ctx, &gen_notify.GetNotificationStatusRequest{
		Id: id,
	})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
