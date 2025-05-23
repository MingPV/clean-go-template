package database

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	sqlDB *sql.DB
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err = db.DB()
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("✅ Database connected")
	return db, nil
}

// Close closes the underlying sql.DB connection pool
func Close() error {
	fmt.Println("🔒 Closing Database connection...")
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}
