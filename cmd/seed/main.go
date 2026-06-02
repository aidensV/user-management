package main

import (
	"log"

	"user-management/internal/config"
	"user-management/internal/database"
	"user-management/internal/models"
	"user-management/internal/utils"
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

	// Seed admin user
	seedAdminUser()
}

func seedAdminUser() {
	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", "admin_user_m@yopmail.com").Count(&count)

	if count > 0 {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	log.Println("Seeding admin user...")

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := &models.User{
		Email:    "admin_user_m@yopmail.com",
		Password: hashedPassword,
		Name:     "Super Admin",
		Role:     "admin",
		IsActive: true,
	}

	if err := database.DB.Create(admin).Error; err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	log.Println("========================================")
	log.Println("Admin user created successfully!")
	log.Println("Email: admin_user_m@yopmail.com")
	log.Println("Password: admin123")
	log.Println("========================================")
}
