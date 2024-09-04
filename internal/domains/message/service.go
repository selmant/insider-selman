package message

import (
	"context"
	"errors"
	"insider/external_services/message_publisher"
	"insider/internal/models"
	"insider/internal/utils"
	"time"
)

type Service interface {
	QueueMessageForSending(ctx context.Context, content, recipientPhone string) error
	FindMessagesBySentStatus(ctx context.Context, sentStatus models.SentStatus) ([]models.Message, error)
	FindAllMessages(ctx context.Context) ([]models.Message, error)
}

type ServiceImpl struct {
	repository       Repository
	messagePublisher message_publisher.MessagePublisher
}

func NewService(repository Repository,
	messagePublisher message_publisher.MessagePublisher) Service {
	return &ServiceImpl{repository: repository,
		messagePublisher: messagePublisher,
	}
}

func (s *ServiceImpl) QueueMessageForSending(ctx context.Context, content, recipientPhone string) error {
	message := models.Message{
		Content:        content,
		RecipientPhone: recipientPhone,
		SentStatus:     models.SentStatusPending,
	}
	_, err := s.repository.Store(ctx, &message)
	return err
}

func (s *ServiceImpl) FindMessagesBySentStatus(ctx context.Context, sentStatus models.SentStatus) ([]models.Message, error) {
	return s.repository.FindBySentStatusWithLimit(ctx, sentStatus, 0)
}

func (s *ServiceImpl) SendQueuedNMessages(ctx context.Context, n int) error {
	messages, err := s.repository.FindBySentStatusWithLimit(ctx, models.SentStatusPending, n)
	if err != nil {
		return err
	}

	var errs error
	for i := range messages {
		message := &messages[i]

		res, err := s.messagePublisher.Publish(ctx, message)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		var updateFields UpdateMessageFields
		updateFields.SentAt = utils.ToPointer(time.Now())
		if res.Message == message_publisher.Accepted {
			updateFields.SentStatus = utils.ToPointer(models.SentStatusSent)
			updateFields.RemoteID = utils.ToPointer(res.MessageID.String())
		} else {
			updateFields.SentStatus = utils.ToPointer(models.SentStatusFailed)
		}

		err = s.repository.Update(ctx, message.ID, updateFields)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
	}

	return errs
}

func (s *ServiceImpl) FindAllMessages(ctx context.Context) ([]models.Message, error) {
	return s.repository.FindAll(ctx)
}
