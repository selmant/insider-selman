package models

import (
	"errors"
	"time"
)

type SentStatus string

type Message struct {
	ID             int        `json:"id" db:"id"`
	Content        string     `json:"content" validate:"required" db:"content"`
	RecipientPhone string     `json:"recipient_phone" validate:"required,e164" db:"recipient_phone"`
	SentStatus     SentStatus `json:"sent_status" validate:"required" db:"sent_status"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	RemoteID       *string    `json:"remote_id" db:"remote_id"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
}

const (
	SentStatusPending SentStatus = "pending"
	SentStatusSent    SentStatus = "sent"
	SentStatusFailed  SentStatus = "failed"
)

func SentStatusFromString(s string) (SentStatus, error) {
	switch s {
	case "pending":
		return SentStatusPending, nil
	case "sent":
		return SentStatusSent, nil
	case "failed":
		return SentStatusFailed, nil
	default:
		return "", errors.New("nknown SentStatus")
	}
}

func (m *Message) IsSent() bool {
	return m.SentStatus == SentStatusSent
}
