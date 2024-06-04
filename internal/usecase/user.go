package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	ur  UserRepository
	sor SensorOwnerRepository
	sr  SensorRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{
		ur:  ur,
		sor: sor,
		sr:  sr,
	}
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if user.Name == "" {
		return nil, ErrInvalidUserName
	}

	err := u.ur.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := u.ur.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.sr.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	so := domain.SensorOwner{
		UserID:   userID,
		SensorID: sensorID,
	}

	err = u.sor.SaveSensorOwner(ctx, so)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	_, err := u.ur.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sos, err := u.sor.GetSensorsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sensors []domain.Sensor

	for _, so := range sos {
		sensor, err := u.sr.GetSensorByID(ctx, so.SensorID)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, *sensor)
	}

	return sensors, nil
}
