package message_publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"insider/internal/models"
	"net/http"
)

type MessagePublisher interface {
	Publish(ctx context.Context, message *models.Message) (Response, error)
}

type MessagePublisherImpl struct {
	client http.Client
}

func NewMessagePublisher() MessagePublisher {
	client := http.Client{Transport: NewBasicAuthTransport()}
	return &MessagePublisherImpl{
		client: client,
	}
}

func (mp *MessagePublisherImpl) Publish(ctx context.Context, message *models.Message) (Response, error) {
	var response Response

	form := Message{
		Content: message.Content,
		To:      message.RecipientPhone,
	}

	payload, err := json.Marshal(form)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(http.MethodPost, "message", bytes.NewReader(payload))
	if err != nil {
		return response, err
	}

	resp, err := mp.client.Do(req.WithContext(ctx))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return response, fmt.Errorf(ErrorPublishMessage, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}
	return response, nil
}
