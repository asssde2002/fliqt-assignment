package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func ensureRolesExist() []models.Role {
	roles := []models.Role{
		{Name: models.Staff},
		{Name: models.Admin},
	}

	err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
	if err != nil {
		log.Printf("Failed to create roles: %v", err)
	}
	return roles
}

func ensureUserExist() []models.User {
	var users []models.User
	for i := 1; i <= 3; i++ {
		username := fmt.Sprintf("test%d", i)
		password := fmt.Sprintf("test%d", i)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to create hash password: %v", err)
		}
		users = append(users, models.User{Username: username, Password: string(hashedPassword)})
	}

	err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error
	if err != nil {
		log.Printf("Failed to create users: %v", err)
	}
	return users
}

func ensureUserRoleExist(roles []models.Role, users []models.User) {
	var userRoles []models.UserRole
	// add staff role to all users
	for i := 0; i < len(users); i++ {
		userID := users[i].ID
		roleID := roles[0].ID
		userRoles = append(userRoles, models.UserRole{UserID: userID, RoleID: roleID})
	}
	// add admin role to the first user
	userRoles = append(userRoles, models.UserRole{UserID: users[0].ID, RoleID: roles[1].ID})
	err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error
	if err != nil {
		log.Printf("Failed to create userroles: %v", err)
	}
}

func ensureUserPunchCard(users []models.User) {
	var punchCards []models.PunchCard
	startDate := 5
	for i := 0; i < 3; i++ {
		userID := users[0].ID
		createdAt := time.Date(2025, time.February, startDate+i, 4, 0, 0, 0, time.UTC)
		pc := models.PunchCard{
			UserID:    userID,
			CreatedAt: createdAt,
			ClockIn:   sql.NullTime{Time: createdAt, Valid: true},
			ClockOut:  sql.NullTime{Time: createdAt.Add(8 * time.Hour), Valid: true},
		}
		punchCards = append(punchCards, pc)
	}

	createdAt := time.Date(2025, time.February, startDate, 4, 0, 0, 0, time.UTC)
	userID := users[1].ID
	pc := models.PunchCard{
		UserID:    userID,
		CreatedAt: createdAt,
		ClockIn:   sql.NullTime{Time: createdAt, Valid: true},
		ClockOut:  sql.NullTime{Time: createdAt.Add(8 * time.Hour), Valid: true},
	}
	punchCards = append(punchCards, pc)
	err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&punchCards).Error
	if err != nil {
		log.Printf("Failed to create punchcards: %v", err)
	}

}

func main() {
	// config.LoadConfig()
	config.LoadTimeZone()

	db.InitDB()
	defer db.CloseDB()
	db.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}, &models.PunchCard{})
	roles := ensureRolesExist()
	users := ensureUserExist()
	ensureUserRoleExist(roles, users)
	ensureUserPunchCard(users)
}
