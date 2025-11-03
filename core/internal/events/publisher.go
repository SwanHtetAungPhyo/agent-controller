package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type Publisher struct {
	nc *nats.Conn
}

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) PublishUserCreated(userID, email, firstName, lastName string) error {
	event := &Event{
		ID:        uuid.New().String(),
		Type:      "user.created",
		Timestamp: time.Now().UTC(),
		Source:    "core-api",
		Data: map[string]interface{}{
			"user_id":    userID,
			"email":      email,
			"name":       fmt.Sprintf("%s %s", firstName, lastName),
			"first_name": firstName,
			"last_name":  lastName,
		},
	}

	return p.publish("user.created", event)
}

func (p *Publisher) PublishUserUpdated(userID, email, firstName, lastName string) error {
	event := &Event{
		ID:        uuid.New().String(),
		Type:      "user.updated",
		Timestamp: time.Now().UTC(),
		Source:    "core-api",
		Data: map[string]interface{}{
			"user_id":    userID,
			"email":      email,
			"name":       fmt.Sprintf("%s %s", firstName, lastName),
			"first_name": firstName,
			"last_name":  lastName,
		},
	}

	return p.publish("user.updated", event)
}

func (p *Publisher) PublishUserDeleted(userID string) error {
	event := &Event{
		ID:        uuid.New().String(),
		Type:      "user.deleted",
		Timestamp: time.Now().UTC(),
		Source:    "core-api",
		Data: map[string]interface{}{
			"user_id": userID,
		},
	}

	return p.publish("user.deleted", event)
}

func (p *Publisher) publish(subject string, event *Event) error {
	if p.nc == nil || !p.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Info().
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Str("subject", subject).
		Msg("Event published successfully")

	return nil
}
