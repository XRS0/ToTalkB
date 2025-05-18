package application

import (
	"fmt"
	"time"

	"event_manager/internal/domain"
)

type EventQueueService struct {
	repository domain.EventQueueRepository
}

func NewEventQueueService(repository domain.EventQueueRepository) *EventQueueService {
	return &EventQueueService{
		repository: repository,
	}
}

func (s *EventQueueService) JoinQueue(eventID string, userID string) error {
	// Проверяем, не находится ли пользователь уже в очереди
	position, err := s.repository.GetPosition(eventID, userID)
	if err == nil && position > 0 {
		return fmt.Errorf("user already in queue")
	}

	// Получаем текущую очередь
	queues, err := s.repository.FindByEventID(eventID)
	if err != nil {
		return err
	}

	// Определяем следующую позицию
	nextPosition := 1
	if len(queues) > 0 {
		nextPosition = len(queues) + 1
	}

	// Создаем новую запись в очереди
	queue := &domain.EventQueue{
		EventID:  eventID,
		UserID:   userID,
		Status:   string(domain.QueueStatusWaiting),
		Position: nextPosition,
	}

	return s.repository.Save(queue)
}

func (s *EventQueueService) LeaveQueue(eventID string, userID string) error {
	// Получаем текущую позицию пользователя
	position, err := s.repository.GetPosition(eventID, userID)
	if err != nil {
		return err
	}

	// Получаем все записи в очереди
	queues, err := s.repository.FindByEventID(eventID)
	if err != nil {
		return err
	}

	// Удаляем запись пользователя
	for _, q := range queues {
		if q.UserID == userID {
			if err := s.repository.Delete(q.ID); err != nil {
				return err
			}
			break
		}
	}

	// Обновляем позиции остальных пользователей
	for _, q := range queues {
		if q.Position > position {
			q.Position--
			if err := s.repository.Update(q); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *EventQueueService) GetQueueStatus(eventID string) ([]*domain.EventQueue, error) {
	return s.repository.FindByEventID(eventID)
}

func (s *EventQueueService) GetUserPosition(eventID string, userID string) (int, error) {
	return s.repository.GetPosition(eventID, userID)
}

func (s *EventQueueService) ProcessNext(eventID string) error {
	// Получаем следующую запись в очереди
	next, err := s.repository.GetNext(eventID)
	if err != nil {
		return err
	}
	if next == nil {
		return fmt.Errorf("no users in queue")
	}

	// Обновляем статус записи
	next.Status = string(domain.QueueStatusActive)
	next.UpdatedAt = time.Now()

	return s.repository.Update(next)
}

func (s *EventQueueService) CloseQueue(eventID string) error {
	// Получаем все записи в очереди
	queues, err := s.repository.FindByEventID(eventID)
	if err != nil {
		return err
	}

	// Обновляем статус всех записей
	for _, queue := range queues {
		if queue.Status == string(domain.QueueStatusWaiting) {
			queue.Status = string(domain.QueueStatusCancelled)
			queue.UpdatedAt = time.Now()
			if err := s.repository.Update(queue); err != nil {
				return err
			}
		}
	}

	return nil
}
