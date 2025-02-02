package models

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"notNull;uniqueIndex;size:255"`
	Password  string    `gorm:"notNull"`
	CreatedAt time.Time `gorm:"notNull"`
}
