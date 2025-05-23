package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// LoadCommon sets common global middleware for the app
func FiberMiddleware(app *fiber.App) {
	app.Use(

		logger.New(), // Logs all requests

		cors.New(cors.Config{
			AllowOrigins: "*", // need to be changed in production
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}),
	)
}
