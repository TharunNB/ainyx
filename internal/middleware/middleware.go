package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prefer a caller-supplied request ID so distributed traces stay intact.
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in Fiber locals so other middleware/handlers can read it.
		c.Locals("requestID", requestID)

		// Set on response so the caller can correlate logs with their request.
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Proceed to the next handler in the chain.
		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		requestID, _ := c.Locals("requestID").(string)

		// Choose log level based on status code for easy alerting.
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		}

		switch {
		case status >= 500:
			log.Error("request completed with server error", fields...)
		case status >= 400:
			log.Warn("request completed with client error", fields...)
		default:
			log.Info("request completed", fields...)
		}

		return err
	}
}
