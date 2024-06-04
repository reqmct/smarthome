package postgres

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		pool,
	}
}

const (
	saveEventQuery              = `INSERT INTO events (timestamp, sensor_serial_number, sensor_id, payload) VALUES ($1, $2, $3, $4);`
	getLastEventBySensorIDQuery = `SELECT * FROM events WHERE sensor_id = $1 AND timestamp = (SELECT MAX(timestamp) FROM events WHERE sensor_id = $1);`
	getEventsByTimeFrameQuery   = `SELECT * FROM events WHERE sensor_id = $1 AND timestamp BETWEEN $2 AND $3`
)

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if event == nil {
		return errors.New("event is nil")
	}

	_, err := r.pool.Exec(ctx, saveEventQuery, event.Timestamp, event.SensorSerialNumber, event.SensorID, event.Payload)
	if err != nil {
		return fmt.Errorf("can't save event: %w", err)
	}

	return nil
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var event domain.Event

	err := r.pool.QueryRow(ctx, getLastEventBySensorIDQuery, id).Scan(&event.Timestamp, &event.SensorSerialNumber, &event.SensorID, &event.Payload)
	if err != nil {
		return nil, fmt.Errorf("can't get last event: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) GetEventsByTimeFrame(ctx context.Context, id int64, start, finish time.Time) ([]domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	rows, err := r.pool.Query(ctx, getEventsByTimeFrameQuery, id, start, finish)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event

	for rows.Next() {
		var event domain.Event
		err = rows.Scan(&event.Timestamp, &event.SensorSerialNumber, &event.SensorID, &event.Payload)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
