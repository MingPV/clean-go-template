package repository

import "github.com/MingPV/clean-go-template/internal/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uint) (*entities.User, error)
	FindAll() ([]entities.User, error)
}
