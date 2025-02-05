package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type PunchCard struct {
	ID        int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64        `json:"user_id" gorm:"notNull;index"`
	ClockIn   sql.NullTime `json:"clock_in"`
	ClockOut  sql.NullTime `json:"clock_out"`
	CreatedAt time.Time    `json:"createdAt" gorm:"notNull"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (pc *PunchCard) BeforeCreate(db *gorm.DB) (err error) {
	currTime := time.Now().UTC()
	pc.CreatedAt = currTime
	pc.ClockIn = sql.NullTime{Time: currTime, Valid: true}
	return
}
