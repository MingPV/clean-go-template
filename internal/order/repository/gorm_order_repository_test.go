package repository_test

import (
	"testing"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/repository"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    repository.OrderRepository
	cleanup func()
}

func (s *OrderRepositoryTestSuite) SetupTest() {
	s.db, s.cleanup = database.SetupTestDB(s.T())
	s.repo = repository.NewGormOrderRepository(s.db)
}

func (s *OrderRepositoryTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TestSave() {
	order := &entities.Order{
		Total: 100.50,
	}

	err := s.repo.Save(order)
	s.NoError(err)
	s.NotZero(order.ID)
}

func (s *OrderRepositoryTestSuite) TestFindByID() {
	// Create an order first
	order := &entities.Order{
		Total: 200.75,
	}
	err := s.repo.Save(order)
	s.NoError(err)

	// Find by ID
	found, err := s.repo.FindByID(int(order.ID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(order.ID, found.ID)
	s.Equal(order.Total, found.Total)
}

func (s *OrderRepositoryTestSuite) TestFindByID_NotFound() {
	_, err := s.repo.FindByID(99999)
	s.Error(err)
}

func (s *OrderRepositoryTestSuite) TestFindAll() {
	// Create multiple orders
	orders := []*entities.Order{
		{Total: 100.0},
		{Total: 200.0},
		{Total: 300.0},
	}

	for _, order := range orders {
		err := s.repo.Save(order)
		s.NoError(err)
	}

	// Find all
	allOrders, err := s.repo.FindAll()
	s.NoError(err)
	s.Len(allOrders, 3)
}

func (s *OrderRepositoryTestSuite) TestFindAll_Empty() {
	allOrders, err := s.repo.FindAll()
	s.NoError(err)
	s.Empty(allOrders)
}

func (s *OrderRepositoryTestSuite) TestPatch() {
	// Create an order first
	order := &entities.Order{
		Total: 150.0,
	}
	err := s.repo.Save(order)
	s.NoError(err)

	// Update order
	updateData := &entities.Order{
		Total: 250.0,
	}
	err = s.repo.Patch(int(order.ID), updateData)
	s.NoError(err)

	// Verify update
	updated, err := s.repo.FindByID(int(order.ID))
	s.NoError(err)
	s.Equal(250.0, updated.Total)
}

func (s *OrderRepositoryTestSuite) TestPatch_NotFound() {
	updateData := &entities.Order{
		Total: 999.0,
	}
	err := s.repo.Patch(99999, updateData)
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *OrderRepositoryTestSuite) TestDelete() {
	// Create an order first
	order := &entities.Order{
		Total: 500.0,
	}
	err := s.repo.Save(order)
	s.NoError(err)

	orderID := int(order.ID)

	// Delete order
	err = s.repo.Delete(orderID)
	s.NoError(err)

	// Verify deletion
	_, err = s.repo.FindByID(orderID)
	s.Error(err)
}

func (s *OrderRepositoryTestSuite) TestDelete_NotFound() {
	err := s.repo.Delete(99999)
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *OrderRepositoryTestSuite) TestSave_MultipleOrders() {
	// Create multiple orders and verify they have different IDs
	order1 := &entities.Order{Total: 100.0}
	order2 := &entities.Order{Total: 200.0}
	order3 := &entities.Order{Total: 300.0}

	err := s.repo.Save(order1)
	s.NoError(err)
	err = s.repo.Save(order2)
	s.NoError(err)
	err = s.repo.Save(order3)
	s.NoError(err)

	s.NotEqual(order1.ID, order2.ID)
	s.NotEqual(order2.ID, order3.ID)
	s.NotEqual(order1.ID, order3.ID)
}
