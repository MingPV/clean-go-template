package repository_test

import (
	"testing"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/user/repository"
	"github.com/MingPV/clean-go-template/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    repository.UserRepository
	cleanup func()
}

func (s *UserRepositoryTestSuite) SetupTest() {
	s.db, s.cleanup = database.SetupTestDB(s.T())
	s.repo = repository.NewGormUserRepository(s.db)
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestSave() {
	user := &entities.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	err := s.repo.Save(user)
	s.NoError(err)
	s.NotEmpty(user.ID)
}

func (s *UserRepositoryTestSuite) TestFindByEmail() {
	// Create a user first
	user := &entities.User{
		Email:    "findbyemail@example.com",
		Password: "password123",
		Name:     "Find By Email User",
	}
	err := s.repo.Save(user)
	s.NoError(err)

	// Find by email
	found, err := s.repo.FindByEmail("findbyemail@example.com")
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.Email, found.Email)
	s.Equal(user.Name, found.Name)
}

func (s *UserRepositoryTestSuite) TestFindByEmail_NotFound() {
	found, err := s.repo.FindByEmail("notfound@example.com")
	s.Error(err)
	s.Nil(found)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *UserRepositoryTestSuite) TestFindByID() {
	// Create a user first
	user := &entities.User{
		Email:    "findbyid@example.com",
		Password: "password123",
		Name:     "Find By ID User",
	}
	err := s.repo.Save(user)
	s.NoError(err)

	// Find by ID
	found, err := s.repo.FindByID(user.ID.String())
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Email, found.Email)
}

func (s *UserRepositoryTestSuite) TestFindByID_NotFound() {
	nonExistentID := uuid.New().String()
	found, err := s.repo.FindByID(nonExistentID)
	s.Error(err)
	s.Nil(found)
}

func (s *UserRepositoryTestSuite) TestFindAll() {
	// Create multiple users
	users := []*entities.User{
		{Email: "user1@example.com", Password: "pass1", Name: "User 1"},
		{Email: "user2@example.com", Password: "pass2", Name: "User 2"},
		{Email: "user3@example.com", Password: "pass3", Name: "User 3"},
	}

	for _, user := range users {
		err := s.repo.Save(user)
		s.NoError(err)
	}

	// Find all
	allUsers, err := s.repo.FindAll()
	s.NoError(err)
	s.Len(allUsers, 3)
}

func (s *UserRepositoryTestSuite) TestFindAll_Empty() {
	allUsers, err := s.repo.FindAll()
	s.NoError(err)
	s.Empty(allUsers)
}

func (s *UserRepositoryTestSuite) TestPatch() {
	// Create a user first
	user := &entities.User{
		Email:    "patch@example.com",
		Password: "password123",
		Name:     "Original Name",
	}
	err := s.repo.Save(user)
	s.NoError(err)

	// Update user
	updateData := &entities.User{
		Name: "Updated Name",
	}
	err = s.repo.Patch(user.ID.String(), updateData)
	s.NoError(err)

	// Verify update
	updated, err := s.repo.FindByID(user.ID.String())
	s.NoError(err)
	s.Equal("Updated Name", updated.Name)
	s.Equal(user.Email, updated.Email) // Email should remain unchanged
}

func (s *UserRepositoryTestSuite) TestPatch_NotFound() {
	nonExistentID := uuid.New().String()
	updateData := &entities.User{
		Name: "Updated Name",
	}
	err := s.repo.Patch(nonExistentID, updateData)
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *UserRepositoryTestSuite) TestDelete() {
	// Create a user first
	user := &entities.User{
		Email:    "delete@example.com",
		Password: "password123",
		Name:     "Delete User",
	}
	err := s.repo.Save(user)
	s.NoError(err)

	// Delete user
	err = s.repo.Delete(user.ID.String())
	s.NoError(err)

	// Verify deletion
	found, err := s.repo.FindByID(user.ID.String())
	s.Error(err)
	s.Nil(found)
}

func (s *UserRepositoryTestSuite) TestDelete_NotFound() {
	nonExistentID := uuid.New().String()
	err := s.repo.Delete(nonExistentID)
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *UserRepositoryTestSuite) TestSave_DuplicateEmail() {
	user1 := &entities.User{
		Email:    "duplicate@example.com",
		Password: "password123",
		Name:     "User 1",
	}
	err := s.repo.Save(user1)
	s.NoError(err)

	// Try to save another user with same email
	user2 := &entities.User{
		Email:    "duplicate@example.com",
		Password: "password456",
		Name:     "User 2",
	}
	err = s.repo.Save(user2)
	s.Error(err) // Should fail due to unique constraint
}
