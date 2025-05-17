package domain

import (
	"github.com/XRS0/ToTalkB/auth/pkg"
)

type EventQueueRepository interface {
	Save(event *Event) error
	FindQueueByEventID(id string) (*EventQueueRepository, error)
	GetFirst() (*pkg.User, error)
}
