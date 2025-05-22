package main

import (
	"log"

	"fmt"
	"os"

	"github.com/MingPV/clean-go-template/entities"
	"github.com/MingPV/clean-go-template/pkg/routes"
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

	// Register routes
	routes.RegisterPublicRoutes(app, db)

	// Not found route
	routes.RegisterNotFoundRoute(app)

	log.Fatal(app.Listen(":8000"))

}
