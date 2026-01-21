package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/MingPV/clean-go-template/internal/entities"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB creates a test database connection and returns a GORM DB instance
// Uses a single test database and cleans tables before/after each test
func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Load .env.dev file - required for tests
	// Try multiple paths to find .env.dev
	envPaths := []string{
		".env.dev",
		"../../.env.dev",
		"../../../.env.dev",
	}

	var envLoaded bool
	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				envLoaded = true
				break
			}
		}
	}

	if !envLoaded {
		t.Fatalf("Failed to load .env.dev file. Please create .env.dev from .env.example. Tried paths: %v", envPaths)
	}

	// Get test database connection details from environment or use defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_TEST_PORT", "5432")
	dbUser := getEnv("DB_TEST_USER", "postgres")
	dbPassword := getEnv("DB_TEST_PASSWORD", "")
	testDBName := getEnv("DB_TEST_NAME", "test")

	// Connect to the test database
	testDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, testDBName,
	)

	db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&entities.User{}, &entities.Order{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Clean up tables before test
	// This ensures each test starts with a clean database
	cleanupTables(db)

	// Return cleanup function that will be called after each test
	cleanup := func() {
		// Clean up tables after test to ensure isolation between tests
		cleanupTables(db)
	}

	return db, cleanup
}

// cleanupTables truncates all test tables to ensure clean state
func cleanupTables(db *gorm.DB) {
	// Truncate tables with CASCADE to handle foreign keys
	// RESTART IDENTITY resets auto-increment counters
	_ = db.Exec("TRUNCATE TABLE users, orders RESTART IDENTITY CASCADE")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
