package postgres

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorRepository struct {
	pool *pgxpool.Pool
}

func NewSensorRepository(pool *pgxpool.Pool) *SensorRepository {
	return &SensorRepository{
		pool: pool,
	}
}

const (
	saveSensorQuery = `INSERT INTO sensors (serial_number, type, current_state, description, is_active, registered_at, last_activity) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`
	getSensorByIDQuery           = `SELECT * FROM sensors WHERE id = $1;`
	getSensorBySerialNumberQuery = `SELECT * FROM sensors WHERE serial_number = $1;`
	getSensorsQuery              = `SELECT * FROM sensors;`
	updateSensorQuery            = `UPDATE sensors SET serial_number = $1, type = $2, current_state = $3, 
                   description = $4, is_active = $5, last_activity = $6 WHERE id = $7;`
)

func (r *SensorRepository) updateSensor(ctx context.Context, sensor *domain.Sensor) error {
	_, err := r.pool.Exec(ctx, updateSensorQuery,
		sensor.SerialNumber,
		sensor.Type,
		sensor.CurrentState,
		sensor.Description,
		sensor.IsActive,
		sensor.LastActivity,
		sensor.ID,
	)
	return err
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if sensor == nil {
		return errors.New("sensor is nil")
	}

	if sensor.ID != 0 {
		return r.updateSensor(ctx, sensor)
	}

	sensor.RegisteredAt = time.Now()

	_, err := r.pool.Exec(ctx, saveSensorQuery,
		sensor.SerialNumber,
		sensor.Type,
		sensor.CurrentState,
		sensor.Description,
		sensor.IsActive,
		sensor.RegisteredAt,
		sensor.LastActivity,
	)
	if err != nil {
		return fmt.Errorf("can't save event: %w", err)
	}

	return nil
}

func sensorMap(row pgx.Row) (*domain.Sensor, error) {
	var sensor domain.Sensor

	err := row.Scan(
		&sensor.ID,
		&sensor.SerialNumber,
		&sensor.Type,
		&sensor.CurrentState,
		&sensor.Description,
		&sensor.IsActive,
		&sensor.RegisteredAt,
		&sensor.LastActivity,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get sensor %w", err)
	}

	return &sensor, nil
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	rows, err := r.pool.Query(ctx, getSensorsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []domain.Sensor

	for rows.Next() {
		s, err := sensorMap(rows)
		if err != nil {
			return nil, fmt.Errorf("can't get sensor %w", err)
		}
		sensors = append(sensors, *s)
	}

	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s, err := sensorMap(r.pool.QueryRow(ctx, getSensorByIDQuery, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecase.ErrSensorNotFound
		}
		return nil, fmt.Errorf("can't get sensor %w", err)
	}

	return s, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s, err := sensorMap(r.pool.QueryRow(ctx, getSensorBySerialNumberQuery, sn))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecase.ErrSensorNotFound
		}
		return nil, fmt.Errorf("can't get sensor %w", err)
	}

	return s, nil
}
