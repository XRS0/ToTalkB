package pkg

import "time"

type Message struct {
	Id        int       `json:"id" db:"id"`
	ChatId    int       `json:"chat_id" db:"chat_id"`
	Sender    string    `json:"sender" db:"sender"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
