package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

func GetRole(rolename models.RoleName) (*models.Role, error) {
	var role models.Role
	if err := db.DB.Where("name = ?", rolename).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
