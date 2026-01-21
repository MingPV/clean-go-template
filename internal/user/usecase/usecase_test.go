package usecase_test

import (
	"os"
	"testing"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/user/repository"
	"github.com/MingPV/clean-go-template/internal/user/usecase"
	"github.com/MingPV/clean-go-template/pkg/apperror"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCaseTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    repository.UserRepository
	service usecase.UserUseCase
	cleanup func()
}

func (s *UserUseCaseTestSuite) SetupTest() {
	s.db, s.cleanup = database.SetupTestDB(s.T())
	s.repo = repository.NewGormUserRepository(s.db)
	s.service = usecase.NewUserService(s.repo)

	// Set JWT_SECRET for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-token-generation")
}

func (s *UserUseCaseTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
	os.Unsetenv("JWT_SECRET")
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}

func (s *UserUseCaseTestSuite) TestRegister() {
	user := &entities.User{
		Email:    "register@example.com",
		Password: "password123",
		Name:     "Register User",
	}

	err := s.service.Register(user)
	s.NoError(err)
	s.NotEmpty(user.ID)

	// Verify password is hashed
	s.NotEqual("password123", user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	s.NoError(err)
}

func (s *UserUseCaseTestSuite) TestRegister_DuplicateEmail() {
	user1 := &entities.User{
		Email:    "duplicate@example.com",
		Password: "password123",
		Name:     "User 1",
	}
	err := s.service.Register(user1)
	s.NoError(err)

	// Try to register with same email
	user2 := &entities.User{
		Email:    "duplicate@example.com",
		Password: "password456",
		Name:     "User 2",
	}
	err = s.service.Register(user2)
	s.Error(err)
	s.Equal(apperror.ErrAlreadyExists, err)
}

func (s *UserUseCaseTestSuite) TestLogin() {
	// Register a user first
	user := &entities.User{
		Email:    "login@example.com",
		Password: "password123",
		Name:     "Login User",
	}
	err := s.service.Register(user)
	s.NoError(err)

	// Login with correct credentials
	token, loggedInUser, err := s.service.Login("login@example.com", "password123")
	s.NoError(err)
	s.NotEmpty(token)
	s.NotNil(loggedInUser)
	s.Equal(user.Email, loggedInUser.Email)
}

func (s *UserUseCaseTestSuite) TestLogin_WrongPassword() {
	// Register a user first
	user := &entities.User{
		Email:    "wrongpass@example.com",
		Password: "password123",
		Name:     "Wrong Pass User",
	}
	err := s.service.Register(user)
	s.NoError(err)

	// Login with wrong password
	token, loggedInUser, err := s.service.Login("wrongpass@example.com", "wrongpassword")
	s.Error(err)
	s.Empty(token)
	s.Nil(loggedInUser)
}

func (s *UserUseCaseTestSuite) TestLogin_UserNotFound() {
	token, loggedInUser, err := s.service.Login("notfound@example.com", "password123")
	s.Error(err)
	s.Empty(token)
	s.Nil(loggedInUser)
}

func (s *UserUseCaseTestSuite) TestFindUserByID() {
	// Register a user first
	user := &entities.User{
		Email:    "findbyid@example.com",
		Password: "password123",
		Name:     "Find By ID User",
	}
	err := s.service.Register(user)
	s.NoError(err)

	// Find by ID
	found, err := s.service.FindUserByID(user.ID.String())
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Email, found.Email)
}

func (s *UserUseCaseTestSuite) TestFindAllUsers() {
	// Register multiple users
	users := []*entities.User{
		{Email: "all1@example.com", Password: "pass1", Name: "User 1"},
		{Email: "all2@example.com", Password: "pass2", Name: "User 2"},
		{Email: "all3@example.com", Password: "pass3", Name: "User 3"},
	}

	for _, user := range users {
		err := s.service.Register(user)
		s.NoError(err)
	}

	// Find all
	allUsers, err := s.service.FindAllUsers()
	s.NoError(err)
	s.Len(allUsers, 3)
}

func (s *UserUseCaseTestSuite) TestPatchUser() {
	// Register a user first
	user := &entities.User{
		Email:    "patch@example.com",
		Password: "password123",
		Name:     "Original Name",
	}
	err := s.service.Register(user)
	s.NoError(err)

	// Update user
	updateData := &entities.User{
		Name: "Updated Name",
	}
	updated, err := s.service.PatchUser(user.ID.String(), updateData)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Name", updated.Name)
	s.Equal(user.Email, updated.Email)
}

func (s *UserUseCaseTestSuite) TestDeleteUser() {
	// Register a user first
	user := &entities.User{
		Email:    "delete@example.com",
		Password: "password123",
		Name:     "Delete User",
	}
	err := s.service.Register(user)
	s.NoError(err)

	// Delete user
	err = s.service.DeleteUser(user.ID.String())
	s.NoError(err)

	// Verify deletion
	found, err := s.service.FindUserByID(user.ID.String())
	s.Error(err)
	s.Nil(found)
}
