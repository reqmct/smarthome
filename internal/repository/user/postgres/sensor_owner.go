package postgres

import (
	"context"
	"fmt"
	"homework/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorOwnerRepository struct {
	pool *pgxpool.Pool
}

func NewSensorOwnerRepository(pool *pgxpool.Pool) *SensorOwnerRepository {
	return &SensorOwnerRepository{
		pool,
	}
}

const (
	saveSensorOwnerQuery    = `INSERT INTO sensors_users (sensor_id, user_id) VALUES ($1, $2);`
	getSensorsByUserIDQuery = `SELECT sensor_id, user_id FROM sensors_users WHERE user_id = $1;`
)

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	_, err := r.pool.Exec(ctx, saveSensorOwnerQuery, sensorOwner.SensorID, sensorOwner.UserID)
	if err != nil {
		return fmt.Errorf("can't save sensor owner: %w", err)
	}

	return nil
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	rows, err := r.pool.Query(ctx, getSensorsByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get sensors: %w", err)
	}
	defer rows.Close()

	var sensors []domain.SensorOwner

	for rows.Next() {
		var s domain.SensorOwner
		err = rows.Scan(&s.SensorID, &s.UserID)
		if err != nil {
			return nil, fmt.Errorf("can't scan sensor: %w", err)
		}
		sensors = append(sensors, s)
	}

	return sensors, nil
}
