package message

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"insider/external_services/message_publisher"
	"insider/internal/models"
	"insider/internal/utils"
	"sync/atomic"
	"time"
)

type Service interface {
	QueueMessageForSending(ctx context.Context, content, recipientPhone string) error
	FindMessagesBySentStatus(ctx context.Context, sentStatus models.SentStatus) ([]models.Message, error)
	FindAllMessages(ctx context.Context) ([]models.Message, error)
	StartMessageSenderJob(ctx context.Context) error
	StopMessageSenderJob(ctx context.Context) error
}

type ServiceImpl struct {
	repository       Repository
	messagePublisher message_publisher.MessagePublisher

	runnerState   atomic.Bool
	jobCancelFunc context.CancelFunc
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

func (s *ServiceImpl) StartMessageSenderJob(_ context.Context) error {
	if s.runnerState.Load() {
		return errors.New("job is already running")
	}
	s.runnerState.Store(true)
	ctx, cancel := context.WithCancel(context.Background())
	s.jobCancelFunc = cancel

	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := s.SendQueuedNMessages(ctx, 2)
				if err != nil {
					log.Error(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (s *ServiceImpl) StopMessageSenderJob(_ context.Context) error {
	if !s.runnerState.Load() {
		return errors.New("job is not running")
	}
	s.runnerState.Store(false)
	if s.jobCancelFunc != nil {
		s.jobCancelFunc()
	}

	return nil
}
