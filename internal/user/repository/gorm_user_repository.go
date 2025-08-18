package repository

import (
	"errors"

	"github.com/MingPV/clean-go-template/internal/entities"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindAll() ([]*entities.User, error) {
	var userValues []entities.User
	if err := r.db.Find(&userValues).Error; err != nil {
		return nil, err
	}
	users := make([]*entities.User, len(userValues))
	for i := range users {
		users[i] = &userValues[i]
	}
	return users, nil
}

func (r *GormUserRepository) Patch(id string, user *entities.User) error {
	result := r.db.Model(&entities.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormUserRepository) Delete(id string) error {
	result := r.db.Delete(&entities.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
