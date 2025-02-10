package services

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

func GetAllRolesMap() (map[models.RoleName]int64, error) {
	roleMap := make(map[models.RoleName]int64)
	roleMapJSON := make(map[string]int64)
	roleMapVal, err := db.RDB.Get(db.CTX, "rolemap").Result()
	if err == redis.Nil {
		var roles []models.Role
		if err := db.DB.Find(&roles).Error; err != nil {
			return nil, err
		}

		for _, role := range roles {
			roleMapJSON[string(role.Name)] = role.ID
		}

		roleMapJSONStr, err := json.Marshal(roleMapJSON)
		if err != nil {
			return nil, err
		}

		if err := db.RDB.Set(db.CTX, "rolemap", roleMapJSONStr, 3600*time.Second).Err(); err != nil {
			return nil, err
		}

		roleMap = make(map[models.RoleName]int64)
		for k, v := range roleMapJSON {
			roleMap[models.RoleName(k)] = v
		}

		return roleMap, nil
	} else {
		fmt.Println("get key")
		err = json.Unmarshal([]byte(roleMapVal), &roleMapJSON)
		if err != nil {
			return nil, err
		}

		for k, v := range roleMapJSON {
			roleMap[models.RoleName(k)] = v
		}

		return roleMap, nil
	}
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
