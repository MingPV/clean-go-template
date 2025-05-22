package routes

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	// Order
	orderAdapters "github.com/MingPV/clean-go-template/adapters/order"
	orderUsecases "github.com/MingPV/clean-go-template/usecases/order"

	// User
	userAdapters "github.com/MingPV/clean-go-template/adapters/user"
	userUsecases "github.com/MingPV/clean-go-template/usecases/user"
)

func RegisterPublicRoutes(app fiber.Router, db *gorm.DB) {

	api := app.Group("/api/v1")

	// === Dependency Wiring ===

	// Order
	orderRepo := orderAdapters.NewGormOrderRepository(db)
	orderService := orderUsecases.NewOrderService(orderRepo)
	orderHandler := orderAdapters.NewHttpOrderHandler(orderService)

	// User
	userRepo := userAdapters.NewGormUserRepository(db)
	userService := userUsecases.NewUserService(userRepo)
	userHandler := userAdapters.NewHttpUserHandler(userService)

	// === Public Routes ===

	// Auth routes (separated from /users)
	authGroup := api.Group("/auth")
	authGroup.Post("/signup", userHandler.Register)
	authGroup.Post("/signin", userHandler.Login)

	// User routes
	userGroup := api.Group("/users")
	userGroup.Get("/", userHandler.FindAllUsers)
	userGroup.Get("/:id", userHandler.FindUserByID)

	// Order routes
	orderGroup := api.Group("/orders")
	orderGroup.Get("/", orderHandler.FindAllOrders)
	orderGroup.Get("/:id", orderHandler.FindOrderByID)
	orderGroup.Post("/", orderHandler.CreateOrder)
	orderGroup.Delete("/:id", orderHandler.DeleteOrder)
	orderGroup.Patch("/:id", orderHandler.PatchOrder)
}
