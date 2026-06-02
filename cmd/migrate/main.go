package main

import (
	"log"

	"user-management/internal/config"
	"user-management/internal/database"
	"user-management/internal/models"
)

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.CloseDB()

	// Run AutoMigrate
	log.Println("Running migrations...")

	err := database.DB.AutoMigrate(
		// &models.User{},
		// &models.UserSession{},
		&models.Permission{},     // ADD THIS LINE
		&models.Role{},           // ADD THIS LINE
		&models.RolePermission{}, // ADD THIS LINE
		&models.UserRole{},       // ADD THIS LINE
		// &models.Menu{},          // ADD THIS LINE
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}
