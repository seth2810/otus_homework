package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	a.logger.Info("scheduler is running...")

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

	ticker := time.NewTicker(cfg.Interval)

	defer ticker.Stop()

	defer a.logger.Info("scheduler is stopping...")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// continue
		}

		now := time.Now()

		events, err := a.storage.ListDayEvents(ctx, now)
		if err != nil {
			return fmt.Errorf("failed to list events: %w", err)
		}

		a.logger.Info(fmt.Sprintf("check events: %s - %v", now, events))

		for _, e := range events {
			if e.NotificationSent {
				continue
			}

			eventStartsAt := e.StartsAt
			eventEndsAt := eventStartsAt.Add(e.Duration)
			eventNotificatiedAt := eventStartsAt.Add(-e.NotifyBefore)

			if !now.After(eventNotificatiedAt) || !now.Before(eventEndsAt) {
				continue
			}

			a.logger.Info(fmt.Sprintf("queue event notification %q: %s", e.ID, eventStartsAt))

			notification := storage.EventNotification{
				ID:       e.ID,
				Title:    e.Title,
				StartsAt: e.StartsAt,
				UserID:   e.OwnerID,
			}

			body, err := json.Marshal(notification)
			if err != nil {
				return fmt.Errorf("failed to marshal notification: %w", err)
			}

			if err := rmq.Publish(ch, queue, body); err != nil {
				return fmt.Errorf("failed to publish message: %w", err)
			}

			e.NotificationSent = true

			if err := a.storage.UpdateEvent(ctx, e.ID, e); err != nil {
				return fmt.Errorf("failed to update event after sent: %w", err)
			}
		}

		yearAgo := now.AddDate(-1, 0, 0)

		a.logger.Info(fmt.Sprintf("remove outdated events: %s", yearAgo))

		// clearing events that occurred more than 1 year ago
		if err := a.storage.RemoveEventsBefore(ctx, yearAgo); err != nil {
			return fmt.Errorf("failed to remove outdated events: %w", err)
		}
	}
}
