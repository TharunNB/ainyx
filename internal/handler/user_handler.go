package handler

import (
	"errors"
	"strconv"
	"user-api/internal/models"
	"user-api/internal/repository"
	"user-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserHandler(svc service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		service:  svc,
		validate: validator.New(),
		logger:   log,
	}
}

// POST /users
// CreateUser handler
// Request: {}
// Response: 201 Created - UserReponse
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("handler.CreateUser: failed to parse the request body",
			zap.String("request_id", requestID(c)),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		h.logger.Error("handler.CreateUser: service error",
			zap.String("request_id", requestID(c)),
			zap.Error(err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GET /users/:id
// GetUserByID
// Response: 200 OK - UserWithAgeResponse
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "id must be a positive integer",
		})
	}

	user, err := h.service.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		h.logger.Error("handler.GetUserByID: service error",
			zap.String("request_id", requestID(c)),
			zap.Int32("user_id", id),
			zap.Error(err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to fetch user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// PUT /users/:id
//
//	Request:  { "name": "Alice Updated", "dob": "1991-03-15" }
//	Response: 200 OK — UserResponse
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "id must be a positive integer",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("handler.UpdateUser: failed to parse request body",
			zap.String("request_id", requestID(c)),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	user, err := h.service.UpdateUser(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		h.logger.Error("handler.UpdateUser: service error",
			zap.String("request_id", requestID(c)),
			zap.Int32("user_id", id),
			zap.Error(err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// DELETE /users/:id
// DeleteUser handles user deletion requests.
//
//	Response: 204 No Content
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "id must be a positive integer",
		})
	}

	err = h.service.DeleteUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		h.logger.Error("handler.DeleteUser: service error",
			zap.String("request_id", requestID(c)),
			zap.Int32("user_id", id),
			zap.Error(err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GET /users
//
//	Query params: page (default 1), limit (default 10, max 100)
//	Response:     200 OK — PaginatedUsersResponse
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	resp, err := h.service.ListUsers(c.Context(), page, limit)
	if err != nil {
		h.logger.Error("handler.ListUsers: service error",
			zap.String("request_id", requestID(c)),
			zap.Error(err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to list users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func parseID(c *fiber.Ctx) (int32, error) {
	raw := c.Params("id")
	id, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return int32(id), nil
}

// requestID reads the request ID stored by the RequestID middleware.
func requestID(c *fiber.Ctx) string {
	id, _ := c.Locals("requestID").(string)
	return id
}
