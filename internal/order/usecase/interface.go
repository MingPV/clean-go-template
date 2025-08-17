package usecase

import "github.com/MingPV/clean-go-template/internal/entities"

type OrderUseCase interface {
	FindAllOrders() ([]*entities.Order, error)
	CreateOrder(order *entities.Order) error
	PatchOrder(id int, order *entities.Order) (*entities.Order, error)
	DeleteOrder(id int) error
	FindOrderByID(id int) (*entities.Order, error)
}
