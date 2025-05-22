package main

import (
	"log"

	"github.com/MingPV/clean-go-template/entities"
	"github.com/MingPV/clean-go-template/pkg/config"
	"github.com/MingPV/clean-go-template/pkg/routes"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config from .env or environment
	cfg := config.LoadConfig()

	// Start Fiber app
	app := fiber.New()

	// Connect to PostgreSQL with config
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}

	// Auto migrate entities
	if err := db.AutoMigrate(&entities.Order{}, &entities.User{}); err != nil {
		log.Fatalf("❌ failed to migrate database: %v", err)
	}

	// Register routes
	routes.RegisterPublicRoutes(app, db)
	routes.RegisterPrivateRoutes(app, db)
	routes.RegisterNotFoundRoute(app)

	// Start server
	log.Printf("🚀 Server is running on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
