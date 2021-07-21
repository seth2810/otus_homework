package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/rmq"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
)

const queue = "notifications"

type App struct {
	logger  Logger
	storage *sqlstorage.Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func New(logger Logger, storage *sqlstorage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Serve(ctx context.Context, cfg *Config) error {
	a.logger.Info("sender is running...")

	conn, err := rmq.Dial("amqp", cfg.RMQ)
	if err != nil {
		return fmt.Errorf("failed to create AMQP connection: %w", err)
	}

	defer conn.Close()

	ch, err := rmq.DeclareQueue(conn, &rmq.QueueDeclareOptions{
		Name:    queue,
		Durable: true,
	})
	if err != nil {
		return fmt.Errorf("failed to prepare AMQP queue: %w", err)
	}

	defer ch.Close()

	deliveryCh, err := rmq.Consume(ctx, ch, queue)
	if err != nil {
		return fmt.Errorf("failed to consume AMQP queue: %w", err)
	}

	var notification *storage.EventNotification

	for m := range deliveryCh {
		notification = &storage.EventNotification{}

		if err := json.Unmarshal(m.Body, notification); err != nil {
			return fmt.Errorf("failed to unmarshal notification: %w", err)
		}

		a.logger.Error(fmt.Sprintf("send event notification: %#v", notification))

		if err := m.Ack(false); err != nil {
			return fmt.Errorf("failed to ack message: %w", err)
		}
	}

	return nil
}
