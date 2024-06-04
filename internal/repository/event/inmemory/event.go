package inmemory

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	mu     sync.Mutex
	events []*domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if event == nil {
		return errors.New("event is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)

	return nil
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var out *domain.Event

	for _, event := range r.events {
		if event.SensorID == id {
			if out == nil {
				out = event
			} else if event.Timestamp.After(out.Timestamp) {
				out = event
			}
		}
		fmt.Println(event.Timestamp)
	}

	if out != nil {
		return out, nil
	}

	return nil, usecase.ErrEventNotFound
}

func (r *EventRepository) GetEventsByTimeFrame(ctx context.Context, id int64, start, finish time.Time) ([]domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var out []domain.Event

	for _, event := range r.events {
		if event.SensorID == id && !(event.Timestamp.Before(start) || event.Timestamp.After(finish)) {
			out = append(out, *event)
		}
	}

	return out, nil
}
