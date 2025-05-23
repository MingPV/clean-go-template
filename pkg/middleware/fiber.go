package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// LoadCommon sets common global middleware for the app
func FiberMiddleware(app *fiber.App) {
	app.Use(

		logger.New(), // Logs all requests

		cors.New(cors.Config{
			AllowOrigins: []string{"*"}, // need to be changed in production
			AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		}),
	)
}
