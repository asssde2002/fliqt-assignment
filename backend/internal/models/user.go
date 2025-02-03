package models

import (
	"os"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" binding:"required" gorm:"notNull;uniqueIndex;size:255"`
	Password  string    `json:"-" binding:"required" gorm:"notNull"`
	CreatedAt time.Time `json:"createdAt" gorm:"notNull"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.CreatedAt = time.Now().UTC()
	return
}

func (user *User) AfterFind(tx *gorm.DB) (err error) {
	timezone := os.Getenv("TIMEZONE")
	userLocation, _ := time.LoadLocation(timezone)
	user.CreatedAt = user.CreatedAt.In(userLocation)
	return
}

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
