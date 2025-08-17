package repository

import (
	"github.com/MingPV/clean-go-template/internal/entities"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) Save(order *entities.Order) error {
	return r.db.Create(&order).Error
}

func (r *GormOrderRepository) FindAll() ([]*entities.Order, error) {
	var orderValues []entities.Order
	if err := r.db.Find(&orderValues).Error; err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, len(orderValues))
	for i := range orderValues {
		orders[i] = &orderValues[i]
	}
	return orders, nil
}

func (r *GormOrderRepository) FindByID(id int) (*entities.Order, error) {
	var order entities.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return &entities.Order{}, err
	}
	return &order, nil
}

func (r *GormOrderRepository) Patch(id int, order *entities.Order) error {
	result := r.db.Model(&entities.Order{}).Where("id = ?", id).Updates(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormOrderRepository) Delete(id int) error {
	result := r.db.Delete(&entities.Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
