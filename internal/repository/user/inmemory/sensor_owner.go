package inmemory

import (
	"context"
	"homework/internal/domain"
	"sync"
)

type SensorOwnerRepository struct {
	mu           sync.Mutex
	sensorOwners map[int64][]domain.SensorOwner
}

func NewSensorOwnerRepository() *SensorOwnerRepository {
	return &SensorOwnerRepository{
		sensorOwners: make(map[int64][]domain.SensorOwner),
	}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sensorOwners[sensorOwner.UserID] = append(r.sensorOwners[sensorOwner.UserID], sensorOwner)

	return nil
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if sensorOwners, ok := r.sensorOwners[userID]; ok {
		return sensorOwners, nil
	}
	return []domain.SensorOwner{}, nil
}
