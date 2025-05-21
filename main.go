package main

import (
	"log"

	"fmt"
	"os"

	adapters "github.com/MingPV/clean-go-template/adapters/order"
	"github.com/MingPV/clean-go-template/entities"
	"github.com/MingPV/clean-go-template/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&entities.Order{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	orderRepo := adapters.NewGormOrderRepository(db)
	orderService := usecases.NewOrderService(orderRepo)
	orderHandler := adapters.NewHttpOrderHandler(orderService)

	app.Get("/orders", orderHandler.FindAllOrders)
	app.Get("/orders/:id", orderHandler.FindOrderByID)
	app.Post("/orders", orderHandler.CreateOrder)
	app.Delete("/orders/:id", orderHandler.DeleteOrder)
	app.Patch("/orders/:id", orderHandler.PatchOrder)

	log.Fatal(app.Listen(":8000"))
}
