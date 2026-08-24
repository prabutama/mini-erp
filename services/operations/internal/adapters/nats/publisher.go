package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/isapr/mini-erp/services/operations/internal/domain"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const streamName = "DOMAIN_EVENTS"

type Publisher struct {
	conn *natsgo.Conn
	js   jetstream.JetStream
}

func NewPublisher(ctx context.Context, url string) (*Publisher, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: streamName, Subjects: []string{"*.>"}}); err != nil {
		conn.Close()
		return nil, err
	}
	return &Publisher{conn: conn, js: js}, nil
}

func (p *Publisher) Publish(ctx context.Context, event domain.EventEnvelope) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	subject := strings.ReplaceAll(event.EventType, "-", "_")
	if subject == "" {
		return fmt.Errorf("event type is required")
	}
	_, err = p.js.Publish(ctx, subject, payload, jetstream.WithMsgID(event.EventID), jetstream.WithRetryWait(100*time.Millisecond), jetstream.WithRetryAttempts(3))
	return err
}

func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
