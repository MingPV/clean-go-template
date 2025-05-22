package usecases

import (
	"errors"

	"github.com/MingPV/clean-go-template/entities"
)

// Tell what OrderService can do
type OrderUseCase interface {
	FindAllOrders() ([]entities.Order, error)
	CreateOrder(order *entities.Order) error
	PatchOrder(id int, order entities.Order) error
	DeleteOrder(id int) error
	FindOrderByID(id int) (entities.Order, error)
}

// OrderService
type OrderService struct {
	repo OrderRepository
}

// Init OrderService function
func NewOrderService(repo OrderRepository) OrderUseCase {
	return &OrderService{repo: repo}
}

// OrderService Methods - 1 create
func (s *OrderService) CreateOrder(order *entities.Order) error {
	if err := s.repo.Save(order); err != nil {
		return err
	}
	return nil
}

// OrderService Methods - 2 find all
func (s *OrderService) FindAllOrders() ([]entities.Order, error) {
	orders, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// OrderService Methods - 3 find by id
func (s *OrderService) FindOrderByID(id int) (entities.Order, error) {

	order, err := s.repo.FindByID(id)
	if err != nil {
		return entities.Order{}, err
	}
	return order, nil
}

// OrderService Methods - 4 patch
func (s *OrderService) PatchOrder(id int, order entities.Order) error {
	if order.Total <= 0 {
		return errors.New("total must be positive")
	}
	if err := s.repo.Patch(id, order); err != nil {
		return err
	}
	return nil
}

// OrderService Methods - 5 delete
func (s *OrderService) DeleteOrder(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
