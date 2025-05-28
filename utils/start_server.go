package utils

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/MingPV/clean-go-template/pkg/config"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func StartRestServer(app *fiber.App, cfg *config.Config) {
	log.Println("Starting REST server on port:", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("REST server error: %v", err)
	}
}

func StartGrpcServer(grpcServer *grpc.Server, cfg *config.Config) {
	log.Println("Starting gRPC server on port:", cfg.GrpcPort)
	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

func WaitForShutdown(cleanups []func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c // Wait for signal
	log.Println("Shutting down...")

	for _, cleanup := range cleanups {
		cleanup()
	}

	log.Println("Shutdown complete.")
}
