package message_publisher

import "github.com/google/uuid"

type ResponseMessage string

const (
	Accepted ResponseMessage = "accepted"
)

type Message struct {
	Content string `json:"content"`
	To      string `json:"to"`
}

type Response struct {
	Message   ResponseMessage `json:"message"`
	MessageID uuid.UUID       `json:"messageId"`
}
