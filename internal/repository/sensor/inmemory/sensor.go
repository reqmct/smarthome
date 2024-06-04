package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"time"
)

var ErrSensorNotFound = errors.New("sensor not found")

type SensorRepository struct {
	muByID     sync.Mutex
	muBySN     sync.Mutex
	score      int64
	senorsByID map[int64]*domain.Sensor
	sensorBySN map[string]*domain.Sensor
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{
		senorsByID: make(map[int64]*domain.Sensor),
		sensorBySN: make(map[string]*domain.Sensor),
	}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if sensor == nil {
		return errors.New("sensor is nil")
	}

	r.muByID.Lock()
	defer r.muByID.Unlock()

	r.muBySN.Lock()
	defer r.muBySN.Unlock()

	r.score++
	sensor.ID = r.score

	sensor.RegisteredAt = time.Now()
	r.senorsByID[sensor.ID] = sensor
	r.sensorBySN[sensor.SerialNumber] = sensor

	return nil
}

func getSensorsBy[T comparable](source map[T]*domain.Sensor) []domain.Sensor {
	var sensors []domain.Sensor

	for _, sensor := range source {
		sensors = append(sensors, *sensor)
	}

	return sensors
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.muBySN.Lock()
	defer r.muBySN.Unlock()
	return getSensorsBy(r.sensorBySN), nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.muByID.Lock()
	defer r.muByID.Unlock()

	if sensor, ok := r.senorsByID[id]; ok {
		return sensor, nil
	}

	return nil, usecase.ErrSensorNotFound
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.muBySN.Lock()
	defer r.muBySN.Unlock()

	if sensor, ok := r.sensorBySN[sn]; ok {
		return sensor, nil
	}

	return nil, usecase.ErrSensorNotFound
}
