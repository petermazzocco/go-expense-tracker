package initializers

import (
	"fmt"
	"go-expense-tracker/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB variable
var DB *gorm.DB

// Initialize Database
func InitDatabase() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or cannot be loaded - using environment variables instead")
	}

	dsn := os.Getenv("DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Assign to global DB
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	// Run migrations
	if err := DB.AutoMigrate(&models.User{}, &models.Expense{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Database is running!")
}
