package message

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"insider/internal/models"
)

type Repository interface {
	Store(ctx context.Context, message *models.Message) (int64, error)
	FindAll(ctx context.Context) ([]models.Message, error)
	Update(ctx context.Context, id int, fields UpdateMessageFields) error
	FindBySentStatusWithLimit(ctx context.Context, status models.SentStatus, limit int) ([]models.Message, error)
}

type RepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Store(ctx context.Context, message *models.Message) (int64, error) {
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO messages (content, recipient_phone, sent_status) VALUES (:content, :recipient_phone, :sent_status)", message)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RepositoryImpl) FindAll(ctx context.Context) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.SelectContext(ctx, &messages, "SELECT * FROM messages")

	return messages, err
}

// FindBySentStatusWithLimit retrieves messages with the specified sent status.
// If the limit is greater than 0, it returns up to the specified number of messages.
// If the limit is 0, it returns all messages with the specified sent status.
func (r *RepositoryImpl) FindBySentStatusWithLimit(ctx context.Context, status models.SentStatus, limit int) ([]models.Message, error) {
	var messages []models.Message
	var err error

	if limit == 0 {
		err = r.db.SelectContext(ctx, &messages, "SELECT * FROM messages WHERE sent_status = ?", status)
	} else {
		err = r.db.SelectContext(ctx, &messages, "SELECT * FROM messages WHERE sent_status = ? LIMIT ?", status, limit)
	}

	return messages, err
}

func (r *RepositoryImpl) Update(ctx context.Context, id int, fields UpdateMessageFields) error {
	query := "UPDATE messages SET "
	values := []interface{}{}

	if fields.Content != nil {
		query += "content = ?, "
		values = append(values, *fields.Content)
	}
	if fields.RecipientPhone != nil {
		query += "recipient_phone = ?, "
		values = append(values, *fields.RecipientPhone)
	}
	if fields.SentStatus != nil {
		query += "sent_status = ?, "
		values = append(values, *fields.SentStatus)
	}
	if fields.SentAt != nil {
		query += "sent_at = ?, "
		values = append(values, *fields.SentAt)
	}

	if len(values) == 0 {
		return errors.New(ErrorMessageNoFields)
	}

	query = query[:len(query)-2] + " WHERE id = ?"
	values = append(values, id)

	resp, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	rowsAffected, err := resp.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(ErrorMessageNotFound)
	}
	return err
}
