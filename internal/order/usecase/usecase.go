package usecase

import (
	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/repository"
)

// OrderService
type OrderService struct {
	repo repository.OrderRepository
}

// Init OrderService function
func NewOrderService(repo repository.OrderRepository) OrderUseCase {
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
func (s *OrderService) FindAllOrders() ([]*entities.Order, error) {
	orders, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// OrderService Methods - 3 find by id
func (s *OrderService) FindOrderByID(id int) (*entities.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return &entities.Order{}, err
	}

	return order, nil
}

// OrderService Methods - 4 patch
func (s *OrderService) PatchOrder(id int, order *entities.Order) (*entities.Order, error) {
	if err := s.repo.Patch(id, order); err != nil {
		return nil, err
	}

	updatedOrder, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

// OrderService Methods - 5 delete
func (s *OrderService) DeleteOrder(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
