package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

var (
	errEventAlreadyExists = errors.New("event already exists")
	errEventNotFound      = errors.New("event not found")
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return errEventAlreadyExists
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return errEventNotFound
	}

	s.events[id] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return errEventNotFound
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return s.listEventsBetween(from, from.AddDate(0, 0, 1))
}

func (s *Storage) ListWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	weekStart := date.AddDate(0, 0, offset)
	from := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, date.Location())

	return s.listEventsBetween(from, from.AddDate(0, 0, 7))
}

func (s *Storage) ListMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	from := time.Date(date.Year(), date.Month(), 0, 0, 0, 0, 0, date.Location())

	return s.listEventsBetween(from, from.AddDate(0, 1, 0))
}

func (s *Storage) listEventsBetween(from, to time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, len(s.events))

	for _, e := range s.events {
		if e.StartsAt.Before(to) && e.StartsAt.After(from) {
			events = append(events, e)
		}
	}

	return events, nil
}
