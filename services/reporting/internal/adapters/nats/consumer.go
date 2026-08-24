package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/isapr/mini-erp/services/reporting/internal/application"
	"github.com/isapr/mini-erp/services/reporting/internal/domain"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const streamName = "DOMAIN_EVENTS"

type Consumer struct {
	conn     *natsgo.Conn
	consume  jetstream.ConsumeContext
	service  *application.ReportingService
	consumer jetstream.Consumer
}

func NewConsumer(ctx context.Context, url string, service *application.ReportingService) (*Consumer, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: streamName, Subjects: []string{"service_order.*"}})
	if err != nil {
		conn.Close()
		return nil, err
	}
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{Durable: "reporting-audit", FilterSubject: "service_order.*", AckPolicy: jetstream.AckExplicitPolicy})
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Consumer{conn: conn, service: service, consumer: consumer}, nil
}

func (c *Consumer) Start() error {
	consume, err := c.consumer.Consume(func(msg jetstream.Msg) {
		var event domain.EventEnvelope
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Printf("invalid domain event: %v", err)
			msg.Ack()
			return
		}
		if _, err := c.service.RecordDomainEvent(context.Background(), event); err != nil {
			log.Printf("record domain event %s failed: %v", event.EventID, err)
			msg.Nak()
			return
		}
		msg.Ack()
	})
	if err != nil {
		return err
	}
	c.consume = consume
	return nil
}

func (c *Consumer) Close() {
	if c.consume != nil {
		c.consume.Stop()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
