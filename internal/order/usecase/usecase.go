package usecase

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/repository"
	"github.com/MingPV/clean-go-template/pkg/redisclient"
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

	// Save to Redis cache
	bytes, _ := json.Marshal(order)
	redisclient.Set("order:"+strconv.FormatUint(uint64(order.ID), 10), string(bytes), time.Minute*10)

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

	// Check if the order is in the cache
	jsonData, err := redisclient.Get("order:" + strconv.Itoa(id))
	if err == nil {
		var order entities.Order
		json.Unmarshal([]byte(jsonData), &order)
		// fmt.Println("Cache hit, returning from cache")
		return &order, nil
	}

	order, err := s.repo.FindByID(id)
	if err != nil {
		return &entities.Order{}, err
	}

	// If not found in the cache, save it to the cache
	// fmt.Println("Cache miss saving to cache")
	bytes, _ := json.Marshal(order)
	redisclient.Set("order:"+strconv.Itoa(id), string(bytes), time.Minute*10)

	return order, nil
}

// OrderService Methods - 4 patch
func (s *OrderService) PatchOrder(id int, order *entities.Order) error {
	if order.Total <= 0 {
		return errors.New("total must be positive")
	}
	if err := s.repo.Patch(id, order); err != nil {
		return err
	}

	// Update cache after patching
	updatedOrder, err := s.repo.FindByID(id)
	if err == nil {
		bytes, _ := json.Marshal(updatedOrder)
		redisclient.Set("order:"+strconv.Itoa(id), string(bytes), time.Minute*10)
	}

	return nil
}

// OrderService Methods - 5 delete
func (s *OrderService) DeleteOrder(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Delete cache after removing from DB
	redisclient.Delete("order:" + strconv.Itoa(id))

	return nil
}
