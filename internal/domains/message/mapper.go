package message

import (
	"insider/internal/models"
	"time"
)

type UpdateMessageFields struct {
	Content        *string            `json:"content" validate:"required" db:"content"`
	RecipientPhone *string            `json:"recipient_phone" validate:"required,e164" db:"recipient_phone"`
	SentStatus     *models.SentStatus `json:"sent_status" validate:"required" db:"sent_status"`
	SentAt         *time.Time         `json:"sent_at" db:"sent_at"`
	RemoteID       *string            `json:"remote_id" db:"remote_id"`
}

type GetMessagesResponse struct {
	Message string           `json:"message"`
	Data    []models.Message `json:"data,omitempty"`
}
