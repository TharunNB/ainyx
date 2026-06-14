package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	db "user-api/db/sqlc"
	"user-api/internal/models"
	"user-api/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const dobLayout = "02-01-2006"

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetUserByID(ctx context.Context, id int32) (models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, limit int) (models.PaginatedUsersResponse, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

// CreateUser implements [UserService].
func (u *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.Dob)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("Invalid dob format: %w", err)
	}

	user, err := u.repo.Create(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob:  pgtype.Date{Time: dob, Valid: true},
	})
	if err != nil {
		u.logger.Error("service.CreateUser: failed to persist user",
			zap.String("name", req.Name),
			zap.Error(err),
		)
		return models.UserResponse{}, err
	}

	u.logger.Info("service.CreateUser: user created",
		zap.Int32("user_id", user.ID),
		zap.String("name", user.Name),
	)
	return toUserResponse(user), nil
}

// DeleteUser implements [UserService].
func (u *userService) DeleteUser(ctx context.Context, id int32) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return repository.ErrNotFound
		}
		u.logger.Error("service.DeleteUser: repository error",
			zap.Int32("user_id", id),
			zap.Error(err),
		)
		return err
	}

	u.logger.Info("service.DeleteUser: user deleted", zap.Int32("user_id", id))
	return nil
}

// GetUserByID implements [UserService].
func (u *userService) GetUserByID(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.UserWithAgeResponse{}, repository.ErrNotFound
		}
		u.logger.Error("service.GetUserById: repo error",
			zap.Int32("user_id", id),
			zap.Error(err),
		)
		return models.UserWithAgeResponse{}, err
	}

	u.logger.Info("service.GetUserById: user fetched",
		zap.Int32("user_id", id),
	)

	return toUserWithAgeResponse(user), nil
}

// ListUsers implements [UserService].
func (u *userService) ListUsers(ctx context.Context, page int, limit int) (models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := u.repo.List(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		u.logger.Error("service.ListUsers: repository error", zap.Error(err))
		return models.PaginatedUsersResponse{}, err
	}

	total, err := u.repo.Count(ctx)
	if err != nil {
		u.logger.Error("service.ListUsers: count error", zap.Error(err))
		return models.PaginatedUsersResponse{}, err
	}

	data := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		data = append(data, toUserWithAgeResponse(u))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	u.logger.Info("service.ListUsers: users listed",
		zap.Int("page", page),
		zap.Int("limit", limit),
		zap.Int64("total", total),
	)

	return models.PaginatedUsersResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser implements [UserService].
func (u *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.Dob)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format: %w", err)
	}

	user, err := u.repo.Update(ctx, db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  pgtype.Date{Time: dob, Valid: true},
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.UserResponse{}, repository.ErrNotFound
		}
		u.logger.Error("service.UpdateUser: repository error",
			zap.Int32("user_id", id),
			zap.Error(err),
		)
		return models.UserResponse{}, err
	}

	u.logger.Info("service.UpdateUser: user updated",
		zap.Int32("user_id", user.ID),
		zap.String("name", user.Name),
	)

	return toUserResponse(user), nil
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func CalculateAge(dob, now time.Time) int {
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return age

}

func toUserResponse(u db.User) models.UserResponse {
	return models.UserResponse{
		ID:   u.ID,
		Name: u.Name,
		Dob:  u.Dob.Time.Format(dobLayout),
	}
}

func toUserWithAgeResponse(u db.User) models.UserWithAgeResponse {
	return models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		Dob:  u.Dob.Time.Format(dobLayout),
		Age:  CalculateAge(u.Dob.Time, time.Now()),
	}
}
