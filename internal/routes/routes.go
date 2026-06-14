// Package routes wires all application routes and global middleware together.
//
// This is the single place where the relationship between URL paths and
// handler functions is declared. Keeping it separate from main.go and the
// handler package means neither needs to know about the other.
package routes

import (
	"user-api/internal/handler"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Register(app *fiber.App, userHandler *handler.UserHandler, log *zap.Logger) {

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	users := app.Group("/users")
	{
		users.Post("/", userHandler.CreateUser)      // POST   /users
		users.Get("/", userHandler.ListUsers)        // GET    /users
		users.Get("/:id", userHandler.GetUserByID)   // GET    /users/:id
		users.Put("/:id", userHandler.UpdateUser)    // PUT    /users/:id
		users.Delete("/:id", userHandler.DeleteUser) // DELETE /users/:id
	}
}
