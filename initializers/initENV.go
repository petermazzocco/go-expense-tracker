package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func InitENV() {
	// Try to load .env file, but don't crash if it's not found
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or cannot be loaded - using environment variables instead")
	}
}
