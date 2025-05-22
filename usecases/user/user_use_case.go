package usecases

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MingPV/clean-go-template/entities"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// What UserService can do
type UserUseCase interface {
	Register(user *entities.User) error
	Login(email, password string) (string, *entities.User, error)
	FindUserByID(id uint) (*entities.User, error)
	FindAllUsers() ([]entities.User, error)
}

// UserService struct
type UserService struct {
	repo UserRepository
}

// Init UserService
func NewUserService(repo UserRepository) UserUseCase {
	return &UserService{repo: repo}
}

// Register user (hash password)
func (s *UserService) Register(user *entities.User) error {
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPwd)

	return s.repo.Save(user)
}

// Login user (check email + password)
func (s *UserService) Login(email string, password string) (string, *entities.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", nil, errors.New("email not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(err)
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

// Get user by ID
func (s *UserService) FindUserByID(id uint) (*entities.User, error) {
	return s.repo.FindByID(id)
}

// Get all users
func (s *UserService) FindAllUsers() ([]entities.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Get user by Email
func (s *UserService) GetUserByEmail(email string) (*entities.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
