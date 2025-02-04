package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
	"log"

	"gorm.io/gorm/clause"
)

func ensureRolesExist() {
	roles := []models.Role{
		{Name: models.Admin},
		{Name: models.Staff},
	}

	err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
	if err != nil {
		log.Printf("Failed to create roles: %v", err)
	}
}

func main() {
	config.LoadConfig()
	db.InitDB()
	defer db.CloseDB()
	db.DB.AutoMigrate(&models.User{}, &models.Role{})
	ensureRolesExist()
}
