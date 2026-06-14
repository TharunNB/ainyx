package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-api/config"
	db "user-api/db/sqlc"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.AppEnv)
	log := logger.Get()
	defer logger.Sync()

	log.Info("starting Ainyx API",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.AppPort),
	)

	ctx := context.Background() // Create a base context at the top of main

	pool := config.ConnectDB(ctx, cfg)
	defer pool.Close()

	log.Info("database connection established")

	if err := config.runMigrations(pool); err != nil {
		log.Fatal("failed to run database migrations", zap.Error(err))
	}
	log.Info("database migrations applied")

	queries := db.New(pool)
	userRepo := repository.NewPostgresUserRepository(queries)
	userSvc := service.NewUserService(userRepo, log)
	userHandler := handler.NewUserHandler(userSvc, log)

	app := fiber.New(fiber.Config{

		ErrorHandler: jsonErrorHandler,

		ServerHeader: "Ainyx",
		AppName:      "Ainyx User API v1.0",
	})

	routes.Register(app, userHandler, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.AppPort)
		log.Info("HTTP server listening", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			log.Error("HTTP server error", zap.Error(err))
		}
	}()

	<-quit
	log.Info("shutdown signal received — draining connections")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Error("forced shutdown after timeout", zap.Error(err))
	}

	log.Info("server stopped cleanly")
}

func jsonErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
