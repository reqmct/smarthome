package usecase

import (
	"context"
	"homework/internal/domain"
	"time"
)

type Event struct {
	er EventRepository
	sr SensorRepository
}

func NewEvent(er EventRepository, sr SensorRepository) *Event {
	return &Event{
		er: er,
		sr: sr,
	}
}

func (e *Event) ReceiveEvent(ctx context.Context, event *domain.Event) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if event.Timestamp.IsZero() {
		return ErrInvalidEventTimestamp
	}

	sensor, err := e.sr.GetSensorBySerialNumber(ctx, event.SensorSerialNumber)
	if err != nil {
		return err
	}

	sensor.LastActivity = time.Now()
	sensor.CurrentState = event.Payload
	event.SensorID = sensor.ID

	if err := e.er.SaveEvent(ctx, event); err != nil {
		return err
	}

	return e.sr.SaveSensor(ctx, sensor)
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return e.er.GetLastEventBySensorID(ctx, id)
}

func (e *Event) GetEventsByTimeFrame(ctx context.Context, id int64, start, finish time.Time) ([]domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return e.er.GetEventsByTimeFrame(ctx, id, start, finish)
}
