package routes

import (
	adapters "github.com/MingPV/clean-go-template/adapters/order"
	"github.com/MingPV/clean-go-template/usecases"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterPublicRoutes(app fiber.Router, db *gorm.DB) {
	orderRepo := adapters.NewGormOrderRepository(db)
	orderService := usecases.NewOrderService(orderRepo)
	orderHandler := adapters.NewHttpOrderHandler(orderService)

	orderGroup := app.Group("/orders")
	orderGroup.Get("/", orderHandler.FindAllOrders)
	orderGroup.Get("/:id", orderHandler.FindOrderByID)
	orderGroup.Post("/", orderHandler.CreateOrder)
	orderGroup.Delete("/:id", orderHandler.DeleteOrder)
	orderGroup.Patch("/:id", orderHandler.PatchOrder)
}
