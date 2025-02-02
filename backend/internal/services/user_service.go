package services

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(newUser *models.User) error {
	var existingUser models.User
	if err := db.DB.Where("username = ?", newUser.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username has already existed")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to create hash password: %v", err)
	}
	newUser.Password = string(hashedPassword)

	if err := db.DB.Create(newUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}
