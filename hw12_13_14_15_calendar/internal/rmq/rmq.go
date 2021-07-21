package rmq

import (
	"context"
	"fmt"

	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type QueueDeclareOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func Dial(protocol string, cfg config.RMQConfig) (*amqp.Connection, error) {
	url := fmt.Sprintf(
		"%s://%s:%s@%s:%d/",
		protocol, cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)

	return amqp.Dial(url)
}

func DeclareQueue(conn *amqp.Connection, opts *QueueDeclareOptions) (*amqp.Channel, error) {
	// Every queue declared gets a default binding to the empty exchange "" which has
	// the type "direct" with the routing key matching the queue's name.
	// With this default binding, it is possible to publish messages that route directly to
	// this queue by publishing to "" with the routing key of the queue name.

	// QueueDeclare("alerts", true, false, false, false, nil)
	// Publish("", "alerts", false, false, Publishing{Body: []byte("...")})

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open AMQP channel: %w", err)
	}

	if _, err := ch.QueueDeclare(opts.Name, opts.Durable, opts.AutoDelete, opts.Exclusive, opts.NoWait, opts.Args); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return ch, nil
}

func Consume(ctx context.Context, ch *amqp.Channel, queue string) (<-chan amqp.Delivery, error) {
	deliveryCh, err := ch.Consume(queue, queue, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to consume channel: %w", err)
	}

	go func() {
		defer ch.Close()

		<-ctx.Done()
	}()

	return deliveryCh, nil
}

func Publish(ch *amqp.Channel, queue string, body []byte) error {
	return ch.Publish("", queue, false, false, amqp.Publishing{
		Body:            body,
		ContentEncoding: "utf-8",
		ContentType:     "application/json",
	})
}
