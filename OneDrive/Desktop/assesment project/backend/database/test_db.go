package database

import (
	"log"

	"shopping-cart/models"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func InitTestDB() {
	var err error
	// Use in-memory SQLite for tests
	DB, err = gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Auto-migrate all models
	DB.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
	)
}

