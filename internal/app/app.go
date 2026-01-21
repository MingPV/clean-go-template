package app

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/MingPV/clean-go-template/internal/entities"
	GrpcOrderHandler "github.com/MingPV/clean-go-template/internal/order/handler/grpc"
	orderRepository "github.com/MingPV/clean-go-template/internal/order/repository"
	orderUseCase "github.com/MingPV/clean-go-template/internal/order/usecase"
	"github.com/MingPV/clean-go-template/pkg/config"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/MingPV/clean-go-template/pkg/middleware"
	"github.com/MingPV/clean-go-template/pkg/routes"
	orderpb "github.com/MingPV/clean-go-template/proto/order"
)

// rest
func SetupRestServer(db *gorm.DB, cfg *config.Config) (*fiber.App, error) {
	app := fiber.New()
	middleware.FiberMiddleware(app)
	// comment out Swagger when testing
	// routes.SwaggerRoute(app)
	routes.RegisterPublicRoutes(app, db)
	routes.RegisterPrivateRoutes(app, db)
	routes.RegisterNotFoundRoute(app)
	return app, nil
}

// grpc
func SetupGrpcServer(db *gorm.DB, cfg *config.Config) (*grpc.Server, error) {
	s := grpc.NewServer()
	orderRepo := orderRepository.NewGormOrderRepository(db)
	orderService := orderUseCase.NewOrderService(orderRepo)

	orderHandler := GrpcOrderHandler.NewGrpcOrderHandler(orderService)
	orderpb.RegisterOrderServiceServer(s, orderHandler)
	return s, nil
}

// dependencies
func SetupDependencies(env string) (*gorm.DB, *config.Config, error) {
	cfg := config.LoadConfig(env)

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		return nil, nil, err
	}

	if env == "test" {
		_ = db.Migrator().DropTable(&entities.Order{}, &entities.User{})
	}
	if err := db.AutoMigrate(&entities.Order{}, &entities.User{}); err != nil {
		return nil, nil, err
	}

	return db, cfg, nil
}
