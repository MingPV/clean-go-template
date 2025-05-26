package routes

import (
	"github.com/MingPV/clean-go-template/internal/user/handler"
	"github.com/MingPV/clean-go-template/internal/user/repository"
	"github.com/MingPV/clean-go-template/internal/user/usecase"
	middleware "github.com/MingPV/clean-go-template/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPrivateRoutes(app fiber.Router, db *gorm.DB) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	userRepo := repository.NewGormUserRepository(db)
	userService := usecase.NewUserService(userRepo)
	userHandler := handler.NewHttpUserHandler(userService)

	route.Get("/me", userHandler.GetUser)

}
