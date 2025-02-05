package services

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"

	"gorm.io/gorm/clause"
)

func GetAllRolesMap() (map[models.RoleName]int64, error) {
	var roles []models.Role
	if err := db.DB.Find(&roles).Error; err != nil {
		return nil, err
	}

	roleMap := make(map[models.RoleName]int64)
	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}
	return roleMap, nil
}

func GetRoleByRoleName(rolename models.RoleName) (*models.Role, error) {
	var role models.Role
	if err := db.DB.Where("name = ?", rolename).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func GetRoleByUserID(userID int64) ([]models.Role, error) {
	var roles []models.Role
	if err := db.DB.
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func PutUserRoles(userID int64, roles []models.RoleName) error {
	roleMap, err := GetAllRolesMap()
	if err != nil {
		return fmt.Errorf("failed to get roles map: %v", err)
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.User{ID: userID}).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove existing roles for user %d: %v", userID, err)
	}

	var newUserRoles []models.UserRole
	for _, roleName := range roles {
		roleID, exists := roleMap[roleName]
		if !exists {
			tx.Rollback()
			return fmt.Errorf("role %s does not exist", roleName)
		}

		newUserRoles = append(newUserRoles, models.UserRole{UserID: userID, RoleID: roleID})
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&newUserRoles).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add roles for user %d: %v", userID, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
