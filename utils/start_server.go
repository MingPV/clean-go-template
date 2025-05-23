package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/MingPV/clean-go-template/pkg/redisclient"
	"github.com/gofiber/fiber/v3"
)

// StartServerWithGracefulShutdown starts the server with graceful shutdown support.
func StartServerWithGracefulShutdown(app *fiber.App, addr string) {
	// Start the server in a goroutine
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Printf("❌ Server error: %v", err)
		}
	}()

	// Create a channel to listen for interrupt/termination signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Wait for signal
	<-ctx.Done()
	log.Println("🛑 Interrupt received. Shutting down server gracefully...")

	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	}

	// Close Redis
	if err := redisclient.CloseRedisClient(); err != nil {
		log.Printf("❌ Failed to close Redis: %v", err)
	} else {
		log.Println("✅ Redis closed successfully.")
	}

	// Close DB connection here
	if err := database.Close(); err != nil {
		log.Printf("❌ Error closing database connection: %v", err)
	} else {
		log.Println("✅ Database connection closed")
	}

	log.Println("👋 Server shutdown complete.")

}

// StartServer starts the server normally without graceful shutdown.
func StartServer(app *fiber.App, addr string) {
	if err := app.Listen(addr); err != nil {
		log.Printf("❌ Server error: %v", err)
	}
}
