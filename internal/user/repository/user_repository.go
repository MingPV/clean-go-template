package repository

import "github.com/MingPV/clean-go-template/internal/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
	FindAll() ([]*entities.User, error)
	Patch(id string, user *entities.User) error
	Delete(id string) error
}
