package domain

import "context"

type EventRepository interface {
	Save(ctx context.Context, event *Event) error
	GetByID(ctx context.Context, id string) (*Event, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Event, error)
}
