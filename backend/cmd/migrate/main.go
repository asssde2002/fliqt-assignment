package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
)

func main() {
	config.LoadConfig()
	db.InitDB()
	defer db.CloseDB()

	db.DB.AutoMigrate(&models.User{})
}
