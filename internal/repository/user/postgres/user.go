package postgres

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

const (
	saveUserQuery    = `INSERT INTO users (name) VALUES ($1);`
	getUserByIDQuery = `SELECT * FROM users WHERE id = $1;`
)

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if user == nil {
		return errors.New("user is nil")
	}

	_, err := r.pool.Exec(ctx, saveUserQuery, user.Name)
	if err != nil {
		return fmt.Errorf("can't save user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var user domain.User

	err := r.pool.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecase.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't get user: %w", err)
	}

	return &user, nil
}
