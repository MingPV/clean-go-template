package main

import (
	"log"

	"github.com/MingPV/clean-go-template/entities"
	"github.com/MingPV/clean-go-template/pkg/config"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/MingPV/clean-go-template/pkg/middleware"
	"github.com/MingPV/clean-go-template/pkg/redisclient"
	"github.com/MingPV/clean-go-template/pkg/routes"
	"github.com/MingPV/clean-go-template/utils"
	"github.com/gofiber/fiber/v2"
)

// @title CleanGO API
// @version 1.0
// @description This is the backend API for CleanGO project.
// @host localhost:8000
// @BasePath /api/v1
func main() {
	// Load config from .env or environment
	cfg := config.LoadConfig()

	// Start Fiber app
	app := fiber.New()

	// Connect to PostgreSQL with config
	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate entities
	if err := db.AutoMigrate(&entities.Order{}, &entities.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Initialize Redis (optional, allow failure in dev mode)
	if err := redisclient.InitRedisClient(cfg.RedisAddress); err != nil {
		log.Printf("redis not available: %v", err)
	}

	middleware.FiberMiddleware(app)

	routes.SwaggerRoute(app)

	// Register routes
	routes.RegisterPublicRoutes(app, db)
	routes.RegisterPrivateRoutes(app, db)
	routes.RegisterNotFoundRoute(app)

	// Start server
	utils.StartServerWithGracefulShutdown(app, ":"+cfg.AppPort)

}
