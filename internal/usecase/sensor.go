package usecase

import (
	"context"
	"errors"
	"homework/internal/domain"
	"log"
	"regexp"
)

type Sensor struct {
	sr SensorRepository
}

func NewSensor(sr SensorRepository) *Sensor {
	return &Sensor{sr: sr}
}

func validateSerialNumber(serialNumber string) bool {
	return regexp.MustCompile(`^(\d\D*){10}$`).MatchString(serialNumber)
}

func (s *Sensor) RegisterSensor(ctx context.Context, sensor *domain.Sensor) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if sensor.Type != domain.SensorTypeContactClosure && sensor.Type != domain.SensorTypeADC {
		return nil, ErrWrongSensorType
	}

	if !validateSerialNumber(sensor.SerialNumber) {
		return nil, ErrWrongSensorSerialNumber
	}

	out, err := s.sr.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
	if !errors.Is(err, ErrSensorNotFound) {
		log.Println(err)
		if err != nil {
			return nil, err
		}
		return out, nil
	}

	err = s.sr.SaveSensor(ctx, sensor)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return s.sr.GetSensors(ctx)
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return s.sr.GetSensorByID(ctx, id)
}
