package routes_test

import (
	"bytes"
	"encoding/json"

	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/MingPV/clean-go-template/internal/app"
)

func setupTestApp(t *testing.T) *fiber.App {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		t.Fatalf("Failed to load .env.test: %v", err)
	}

	db, _, cfg, err := app.SetupDependencies("test")
	if err != nil {
		t.Fatalf("failed to setup dependencies: %v", err)
	}

	restApp, err := app.SetupRestServer(db, cfg)
	if err != nil {
		t.Fatalf("failed to setup REST server: %v", err)
	}

	return restApp
}

func TestPublicRoutes(t *testing.T) {
	app := setupTestApp(t)

	// === USERS ===
	t.Run("GET /api/v1/users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("GET /api/v1/users/:id (not found)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/9a176ca5-f3e0-4994-869c-fac0e8c9d5dc", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	// === AUTH ===
	t.Run("POST /api/v1/auth/signup", func(t *testing.T) {
		body := map[string]string{
			"email":    "testuser@example.com",
			"password": "securepassword123",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusCreated)
	})

	t.Run("POST /api/v1/auth/signin", func(t *testing.T) {
		body := map[string]string{
			"email":    "testuser@example.com",
			"password": "securepassword123",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusUnauthorized)
	})

	// === ORDERS ===
	t.Run("GET /api/v1/orders", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/orders", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("GET /api/v1/orders/:id (not found)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/orders/999", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("POST /api/v1/orders", func(t *testing.T) {
		body := map[string]interface{}{
			"total": 300,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusCreated)
	})

	t.Run("PATCH /api/v1/orders/:id", func(t *testing.T) {
		body := map[string]interface{}{
			"total": 3001,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PATCH", "/api/v1/orders/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 500)
	})

	t.Run("DELETE /api/v1/orders/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/orders/1", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 500)
	})
}
