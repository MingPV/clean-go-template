package routes

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterNotFoundRoute sets a custom handler for 404 Not Found
func RegisterNotFoundRoute(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist.",
		})
	})
}
