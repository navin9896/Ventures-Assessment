package database

import (
	"fmt"
	"log"
	"os"

	"shopping-cart/models"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDB() {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	var err error
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "shopping_cart")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate all models
	DB.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
	)

	log.Println("Database connected and migrated successfully")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func SeedData() {
	// Check if items already exist
	var count int
	DB.Model(&models.Item{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return
	}

	// Seed sample items
	items := []models.Item{
		{Name: "Laptop", Status: "active"},
		{Name: "Mouse", Status: "active"},
		{Name: "Keyboard", Status: "active"},
		{Name: "Monitor", Status: "active"},
		{Name: "Headphones", Status: "active"},
	}

	for _, item := range items {
		DB.Create(&item)
	}

	log.Println("Database seeded with sample items")
}

