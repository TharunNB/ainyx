package repository

import (
	"context"
	"errors"
	"fmt"
	db "user-api/db/sqlc"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("resouce not found")

type UserRepository interface {
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	Update(ctx context.Context, params db.UpdateUserParams) (db.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, params db.ListUsersParams) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
}

type postgresUserRepository struct {
	queries db.Querier
}

// Count implements [UserRepository].
func (p *postgresUserRepository) Count(ctx context.Context) (int64, error) {
	count, err := p.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("repository.Count: %w", err)
	}
	return count, nil
}

// Create implements [UserRepository].
func (p *postgresUserRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	user, err := p.queries.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, fmt.Errorf("repository.Create: %w", err)
	}
	return user, nil
}

// Delete implements [UserRepository].
func (p *postgresUserRepository) Delete(ctx context.Context, id int32) error {
	if _, err := p.GetByID(ctx, id); err != nil {
		return err
	}

	if err := p.queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}
	return nil
}

// GetByID implements [UserRepository].
func (p *postgresUserRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	user, err := p.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, fmt.Errorf("repository.GetByID: %w", err)
	}
	return user, nil
}

// List implements [UserRepository].
func (p *postgresUserRepository) List(ctx context.Context, params db.ListUsersParams) ([]db.User, error) {
	users, err := p.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("repository.List: %w", err)
	}
	return users, nil
}

// Update implements [UserRepository].
func (p *postgresUserRepository) Update(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	user, err := p.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, fmt.Errorf("repository.Update: %w", err)
	}
	return user, nil
}

func NewPostgresUserRepository(queries db.Querier) UserRepository {
	return &postgresUserRepository{queries: queries}
}
