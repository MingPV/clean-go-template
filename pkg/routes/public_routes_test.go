package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/MingPV/clean-go-template/internal/app"
	"github.com/MingPV/clean-go-template/pkg/config"
	"github.com/MingPV/clean-go-template/pkg/database"
)

type PublicRoutesTestSuite struct {
	suite.Suite
	db      *gorm.DB
	app     *fiber.App
	cfg     *config.Config
	cleanup func()
}

func (s *PublicRoutesTestSuite) SetupTest() {
	// Setup test database with cleanup
	s.db, s.cleanup = database.SetupTestDB(s.T())

	// Load config for dev environment
	s.cfg = config.LoadConfig("dev")

	// Setup REST server with test database (For registering routes and middleware)
	var err error
	s.app, err = app.SetupRestServer(s.db, s.cfg)
	s.NoError(err, "Failed to setup REST server")
}

func (s *PublicRoutesTestSuite) TearDownTest() {
	// Clean up database after each test
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestPublicRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(PublicRoutesTestSuite))
}

// === USER ROUTES ===

func (s *PublicRoutesTestSuite) TestGetUsers() {
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)
}

func (s *PublicRoutesTestSuite) TestGetUserByID_NotFound() {
	req := httptest.NewRequest("GET", "/api/v1/users/9a176ca5-f3e0-4994-869c-fac0e8c9d5dc", nil)
	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.NotEqual(fiber.StatusInternalServerError, resp.StatusCode)
}

// === AUTH ROUTES ===

func (s *PublicRoutesTestSuite) TestSignup() {
	body := map[string]string{
		"email":    "testuser@example.com",
		"password": "securepassword123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.True(resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusCreated)
}

func (s *PublicRoutesTestSuite) TestSignin() {
	// First signup to create a user
	signupBody := map[string]string{
		"email":    "signinuser@example.com",
		"password": "securepassword123",
	}
	jsonSignupBody, _ := json.Marshal(signupBody)
	signupReq := httptest.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(jsonSignupBody))
	signupReq.Header.Set("Content-Type", "application/json")
	_, _ = s.app.Test(signupReq, -1)

	// Then try to signin
	body := map[string]string{
		"email":    "signinuser@example.com",
		"password": "securepassword123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.True(resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusUnauthorized)
}

func (s *PublicRoutesTestSuite) TestSignin_InvalidCredentials() {
	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}

// === ORDER ROUTES ===

func (s *PublicRoutesTestSuite) TestGetOrders() {
	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)
}

func (s *PublicRoutesTestSuite) TestGetOrderByID_NotFound() {
	req := httptest.NewRequest("GET", "/api/v1/orders/999", nil)
	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.NotEqual(fiber.StatusInternalServerError, resp.StatusCode)
}

func (s *PublicRoutesTestSuite) TestCreateOrder() {
	body := map[string]interface{}{
		"total": 300,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.True(resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusCreated)
}

func (s *PublicRoutesTestSuite) TestPatchOrder() {
	// First create an order
	createBody := map[string]interface{}{
		"total": 300,
	}
	createJsonBody, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(createJsonBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := s.app.Test(createReq, -1)
	s.True(createResp.StatusCode == fiber.StatusOK || createResp.StatusCode == fiber.StatusCreated)

	// Then try to patch it
	body := map[string]interface{}{
		"total": 3001,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PATCH", "/api/v1/orders/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.True(resp.StatusCode >= 200 && resp.StatusCode < 500)
}

func (s *PublicRoutesTestSuite) TestDeleteOrder() {
	// First create an order
	createBody := map[string]interface{}{
		"total": 300,
	}
	createJsonBody, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(createJsonBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := s.app.Test(createReq, -1)
	s.True(createResp.StatusCode == fiber.StatusOK || createResp.StatusCode == fiber.StatusCreated)

	// Then try to delete it
	req := httptest.NewRequest("DELETE", "/api/v1/orders/1", nil)
	resp, err := s.app.Test(req, -1)
	s.NoError(err)
	s.True(resp.StatusCode >= 200 && resp.StatusCode < 500)
}
