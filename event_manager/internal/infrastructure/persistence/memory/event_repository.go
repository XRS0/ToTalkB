package memory

import (
	"context"
	"fmt"
	"sync"

	"event_manager/internal/domain"
)

type InMemoryEventRepository struct {
	events map[string]*domain.Event
	mu     sync.RWMutex
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events: make(map[string]*domain.Event),
	}
}

func (r *InMemoryEventRepository) Save(ctx context.Context, event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if event.ID == "" {
		return fmt.Errorf("event ID is required")
	}

	r.events[event.ID] = event
	return nil
}

func (r *InMemoryEventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, fmt.Errorf("event not found: %s", id)
	}

	return event, nil
}

func (r *InMemoryEventRepository) Update(ctx context.Context, event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if event.ID == "" {
		return fmt.Errorf("event ID is required")
	}

	if _, exists := r.events[event.ID]; !exists {
		return fmt.Errorf("event not found: %s", event.ID)
	}

	r.events[event.ID] = event
	return nil
}

func (r *InMemoryEventRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return fmt.Errorf("event not found: %s", id)
	}

	delete(r.events, id)
	return nil
}

func (r *InMemoryEventRepository) List(ctx context.Context) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]*domain.Event, 0, len(r.events))
	for _, event := range r.events {
		events = append(events, event)
	}

	return events, nil
}
