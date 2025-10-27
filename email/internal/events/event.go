package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type EventService struct {
	nc            *nats.Conn
	subscriptions []*nats.Subscription
}

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
}

type EventHandler func(event *Event) error

func NewEventService(nc *nats.Conn) *EventService {
	return &EventService{
		nc:            nc,
		subscriptions: make([]*nats.Subscription, 0),
	}
}

func (es *EventService) Start(ctx context.Context) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	log.Println("Event service started")
	return nil
}

func (es *EventService) Stop(ctx context.Context) error {
	log.Println("Stopping event service...")

	for _, sub := range es.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("Error unsubscribing: %v", err)
		}
	}

	es.subscriptions = nil
	log.Println("Event service stopped")
	return nil
}

func (es *EventService) Publish(subject string, event *Event) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := es.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event: subject=%s, type=%s, id=%s", subject, event.Type, event.ID)
	return nil
}

func (es *EventService) PublishAsync(subject string, event *Event) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := es.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	go func() {
		if err := es.nc.FlushTimeout(5 * time.Second); err != nil {
			log.Printf("Failed to flush NATS: %v", err)
		}
	}()

	return nil
}

func (es *EventService) Subscribe(subject string, handler EventHandler) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	sub, err := es.nc.Subscribe(subject, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		if err := handler(&event); err != nil {
			log.Printf("Error handling event: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	es.subscriptions = append(es.subscriptions, sub)
	log.Printf("Subscribed to subject: %s", subject)
	return nil
}

func (es *EventService) QueueSubscribe(subject, queue string, handler EventHandler) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	sub, err := es.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		if err := handler(&event); err != nil {
			log.Printf("Error handling event: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to queue subscribe to subject %s: %w", subject, err)
	}

	es.subscriptions = append(es.subscriptions, sub)
	log.Printf("Queue subscribed to subject: %s (queue: %s)", subject, queue)
	return nil
}

func (es *EventService) Request(subject string, event *Event, timeout time.Duration) (*Event, error) {
	if es.nc == nil || !es.nc.IsConnected() {
		return nil, fmt.Errorf("NATS connection is not available")
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	msg, err := es.nc.Request(subject, data, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var response Event
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (es *EventService) SubscribeWithReply(subject string, handler func(*Event) (*Event, error)) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	sub, err := es.nc.Subscribe(subject, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		response, err := handler(&event)
		if err != nil {
			log.Printf("Error handling request: %v", err)
			errorResponse := &Event{
				Type:      "error",
				Timestamp: time.Now().UTC(),
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			}
			responseData, _ := json.Marshal(errorResponse)
			msg.Respond(responseData)
			return
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return
		}

		if err := msg.Respond(responseData); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe with reply to subject %s: %w", subject, err)
	}

	es.subscriptions = append(es.subscriptions, sub)
	log.Printf("Subscribed with reply to subject: %s", subject)
	return nil
}

func (es *EventService) GetStats() nats.Statistics {
	if es.nc == nil {
		return nats.Statistics{}
	}
	return es.nc.Stats()
}

func (es *EventService) IsConnected() bool {
	return es.nc != nil && es.nc.IsConnected()
}
