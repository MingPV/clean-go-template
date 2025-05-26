package usecase

import "github.com/MingPV/clean-go-template/internal/entities"

// What UserService can do
type UserUseCase interface {
	Register(user *entities.User) error
	Login(email, password string) (string, *entities.User, error)
	FindUserByID(id uint) (*entities.User, error)
	FindAllUsers() ([]entities.User, error)
}
