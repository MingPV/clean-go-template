package usecase_test

import (
	"testing"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/repository"
	"github.com/MingPV/clean-go-template/internal/order/usecase"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OrderUseCaseTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    repository.OrderRepository
	service usecase.OrderUseCase
	cleanup func()
}

func (s *OrderUseCaseTestSuite) SetupTest() {
	s.db, s.cleanup = database.SetupTestDB(s.T())
	s.repo = repository.NewGormOrderRepository(s.db)
	s.service = usecase.NewOrderService(s.repo)
}

func (s *OrderUseCaseTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestOrderUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUseCaseTestSuite))
}

func (s *OrderUseCaseTestSuite) TestCreateOrder() {
	order := &entities.Order{
		Total: 150.50,
	}

	err := s.service.CreateOrder(order)
	s.NoError(err)
	s.NotZero(order.ID)
}

func (s *OrderUseCaseTestSuite) TestFindAllOrders() {
	// Create multiple orders
	orders := []*entities.Order{
		{Total: 100.0},
		{Total: 200.0},
		{Total: 300.0},
	}

	for _, order := range orders {
		err := s.service.CreateOrder(order)
		s.NoError(err)
	}

	// Find all
	allOrders, err := s.service.FindAllOrders()
	s.NoError(err)
	s.Len(allOrders, 3)
}

func (s *OrderUseCaseTestSuite) TestFindAllOrders_Empty() {
	allOrders, err := s.service.FindAllOrders()
	s.NoError(err)
	s.Empty(allOrders)
}

func (s *OrderUseCaseTestSuite) TestFindOrderByID() {
	// Create an order first
	order := &entities.Order{
		Total: 250.75,
	}
	err := s.service.CreateOrder(order)
	s.NoError(err)

	// Find by ID
	found, err := s.service.FindOrderByID(int(order.ID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(order.ID, found.ID)
	s.Equal(order.Total, found.Total)
}

func (s *OrderUseCaseTestSuite) TestFindOrderByID_NotFound() {
	_, err := s.service.FindOrderByID(99999)
	s.Error(err)
}

func (s *OrderUseCaseTestSuite) TestPatchOrder() {
	// Create an order first
	order := &entities.Order{
		Total: 100.0,
	}
	err := s.service.CreateOrder(order)
	s.NoError(err)

	orderID := int(order.ID)

	// Update order
	updateData := &entities.Order{
		Total: 500.0,
	}
	updated, err := s.service.PatchOrder(orderID, updateData)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(500.0, updated.Total)
	s.Equal(orderID, int(updated.ID))
}

func (s *OrderUseCaseTestSuite) TestPatchOrder_NotFound() {
	updateData := &entities.Order{
		Total: 999.0,
	}
	updated, err := s.service.PatchOrder(99999, updateData)
	s.Error(err)
	s.Nil(updated)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *OrderUseCaseTestSuite) TestDeleteOrder() {
	// Create an order first
	order := &entities.Order{
		Total: 600.0,
	}
	err := s.service.CreateOrder(order)
	s.NoError(err)

	orderID := int(order.ID)

	// Delete order
	err = s.service.DeleteOrder(orderID)
	s.NoError(err)

	// Verify deletion
	_, err = s.service.FindOrderByID(orderID)
	s.Error(err)
}

func (s *OrderUseCaseTestSuite) TestDeleteOrder_NotFound() {
	err := s.service.DeleteOrder(99999)
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *OrderUseCaseTestSuite) TestCreateOrder_ZeroTotal() {
	order := &entities.Order{
		Total: 0.0,
	}

	err := s.service.CreateOrder(order)
	s.NoError(err)
	s.NotZero(order.ID)
	s.Equal(0.0, order.Total)
}

func (s *OrderUseCaseTestSuite) TestCreateOrder_LargeTotal() {
	order := &entities.Order{
		Total: 999999.99,
	}

	err := s.service.CreateOrder(order)
	s.NoError(err)
	s.NotZero(order.ID)
	s.Equal(999999.99, order.Total)
}
