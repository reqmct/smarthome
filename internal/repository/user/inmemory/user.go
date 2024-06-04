package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	mu    sync.Mutex
	users map[int64]*domain.User
	score int64
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int64]*domain.User),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if user == nil {
		return errors.New("user is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.score++
	user.ID = r.score
	r.users[user.ID] = user

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if user, ok := r.users[id]; ok {
		return user, ctx.Err()
	}
	return nil, usecase.ErrUserNotFound
}
