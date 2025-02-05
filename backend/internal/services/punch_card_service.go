package services

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"
	"time"
)

func CheckPunchCardExist(today string, userID int64) (*models.PunchCard, bool) {
	var existingPunchCard models.PunchCard
	if err := db.DB.Where("user_id = ? AND DATE(created_at) = ?", userID, today).First(&existingPunchCard).Error; err == nil {
		return &existingPunchCard, true
	}
	return nil, false
}

func CreatePunchCard(userID int64) error {
	punchCard := models.PunchCard{
		UserID: userID,
	}
	if err := db.DB.Create(&punchCard).Error; err != nil {
		return fmt.Errorf("failed to create punch card: %v", err)
	}
	return nil
}

func UpdatePunchCard(punchCardID int64) error {
	if err := db.DB.Model(&models.PunchCard{}).Where("id = ?", punchCardID).Update("clock_out", time.Now().UTC()).Error; err != nil {
		return fmt.Errorf("failed to update punch card: %v", err)
	}
	return nil
}

func GetPunchCard(userID int64) ([]models.PunchCard, error) {
	var punchCards []models.PunchCard
	if err := db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&punchCards).Error; err != nil {
		return nil, err
	}
	return punchCards, nil
}
