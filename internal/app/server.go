package app

import (
	"log"

	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/MingPV/clean-go-template/utils"
)

func Start() {

	// Setup dependencies: database and configuration
	db, cfg, err := SetupDependencies("dev")
	if err != nil {
		log.Fatalf("❌ Failed to setup dependencies: %v", err)
	}

	// Setup REST server
	restApp, err := SetupRestServer(db, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to setup REST server: %v", err)
	}

	// Setup gRPC server
	grpcServer, err := SetupGrpcServer(db, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to setup gRPC server: %v", err)
	}

	// Start REST and gRPC servers
	go utils.StartRestServer(restApp, cfg)
	go utils.StartGrpcServer(grpcServer, cfg)

	// Graceful shutdown listener
	utils.WaitForShutdown([]func(){
		func() {
			log.Println("Shutting down REST server...")
			if err := restApp.Shutdown(); err != nil {
				log.Printf("Error shutting down REST server: %v", err)
			}
		},
		func() {
			log.Println("Shutting down gRPC server...")
			grpcServer.GracefulStop()
		},
		func() {
			if err := database.Close(); err != nil {
				log.Printf("Error closing DB: %v", err)
			}
		},
	})

}
