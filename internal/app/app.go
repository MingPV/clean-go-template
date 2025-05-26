package app

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/pkg/config"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/MingPV/clean-go-template/pkg/middleware"
	"github.com/MingPV/clean-go-template/pkg/redisclient"
	"github.com/MingPV/clean-go-template/pkg/routes"
)

func SetupApp() (*fiber.App, string) {
	// Load config
	cfg := config.LoadConfig()

	// Setup Fiber
	app := fiber.New()

	// Middleware
	middleware.FiberMiddleware(app)

	// Database
	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate entities
	if err := db.AutoMigrate(&entities.Order{}, &entities.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Redis
	if err := redisclient.InitRedisClient(cfg.RedisAddress); err != nil {
		log.Printf("redis not available: %v", err)
	}

	// Swagger
	routes.SwaggerRoute(app)

	// Routes
	routes.RegisterPublicRoutes(app, db)
	routes.RegisterPrivateRoutes(app, db)
	routes.RegisterNotFoundRoute(app)

	return app, cfg.AppPort
}
