package message

import "time"

type Message struct {
	ID             int        `json:"id" db:"id"`
	Content        string     `json:"content" validate:"required" db:"content"`
	RecipientPhone string     `json:"recipient_phone" validate:"required,e164" db:"recipient_phone"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
}
