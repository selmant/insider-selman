package message

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository interface {
	Store(message *Message) error
	FindAll() ([]*Message, error)
	FindByID(id int) (*Message, error)
	Delete(id int) error
	FindNotSent() ([]*Message, error)
}

type RepositoryImpl struct {
	db sqlx.DB
}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (r *RepositoryImpl) Store(ctx context.Context, message *Message) (int64, error) {
	res, err := r.db.NamedExec("INSERT INTO messages (content, recipient_phone) VALUES (:content, :recipient_phone)", message)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RepositoryImpl) FindAll(ctx context.Context) ([]Message, error) {
	var messages []Message
	err := r.db.SelectContext(ctx, &messages, "SELECT * FROM messages")

	return messages, err
}
