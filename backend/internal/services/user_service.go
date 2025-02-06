package services

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(authInput *models.AuthInput) error {
	if _, err := FetchUserByUsername(authInput.Username); err == nil {
		return fmt.Errorf("username has already existed")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to create hash password: %v", err)
	}

	staffRole, err := GetRoleByRoleName(models.Staff)
	if err != nil {
		return fmt.Errorf("cannot get staff role instance")
	}

	newUser := models.User{
		Username: authInput.Username,
		Password: string(hashedPassword),
		Roles:    []models.Role{*staffRole},
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func FetchUserByID(id int64) (*models.User, error) {
	var user models.User
	err := db.DB.
		Joins("LEFT JOIN user_roles ON user_roles.user_id = users.id").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id").
		Where("users.id = ?", id).
		Preload("Roles").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func VerifyUser(authInput *models.AuthInput) (*models.User, error) {
	existingUser, err := FetchUserByUsername(authInput.Username)
	if err != nil {
		return nil, fmt.Errorf("username does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(authInput.Password)); err != nil {
		return nil, fmt.Errorf("failed to get hash password: %v", err)
	}

	if !existingUser.IsActive {
		return nil, fmt.Errorf("user is not acitve: %v", err)
	}

	return existingUser, nil
}
