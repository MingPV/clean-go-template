package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/swagger"
)

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(a *fiber.App) {

	a.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
	}))

}
