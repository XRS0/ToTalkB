package pkg

type Message struct {
	ID       int    `json:"id"`
	ChatID   int    `json:"chat_id"`
	SenderID int    `json:"sender_id"`
	Content  string `json:"content"`
}
