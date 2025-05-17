package pkg

type Message struct {
	Id         string `json:"id"`
	ChatId     string `json:"chat_id"`
	SenderId   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
}
