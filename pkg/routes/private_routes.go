package routes

import (
	userAdapters "github.com/MingPV/clean-go-template/adapters/user"
	middleware "github.com/MingPV/clean-go-template/pkg/middleware"
	userUsecases "github.com/MingPV/clean-go-template/usecases/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPrivateRoutes(app fiber.Router, db *gorm.DB) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	userRepo := userAdapters.NewGormUserRepository(db)
	userService := userUsecases.NewUserService(userRepo)
	userHandler := userAdapters.NewHttpUserHandler(userService)

	route.Get("/me", userHandler.GetUser)

}
