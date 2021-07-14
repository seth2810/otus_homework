package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	ownerID := uuid.New()

	return a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title, OwnerID: ownerID.String()})
}

func (a *App) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListDayEvents(ctx, date)
}

func (a *App) ListWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListWeekEvents(ctx, date)
}

func (a *App) ListMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListMonthEvents(ctx, date)
}
