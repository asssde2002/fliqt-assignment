package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" binding:"required" gorm:"notNull;uniqueIndex;size:255"`
	Password  string    `json:"-" binding:"required" gorm:"notNull"`
	CreatedAt time.Time `json:"created_at" gorm:"notNull"`
	IsActive  bool      `json:"is_active" gorm:"notNull"`
	Roles     []Role    `json:"roles" gorm:"many2many:user_roles"`
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	user.CreatedAt = time.Now().UTC()
	user.IsActive = true
	return
}

func (user *User) AfterFind(db *gorm.DB) (err error) {
	user.CreatedAt = user.CreatedAt.Local()
	return
}

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	IsActive  bool       `json:"is_active"`
	Roles     []RoleName `json:"roles"`
}
