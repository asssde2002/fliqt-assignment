package models

import (
	"database/sql"
	"os"
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

func (pc *PunchCard) AfterFind(db *gorm.DB) (err error) {
	timezone := os.Getenv("TIMEZONE")
	userLocation, _ := time.LoadLocation(timezone)
	pc.CreatedAt = pc.CreatedAt.In(userLocation)
	if pc.ClockIn.Valid {
		pc.ClockIn.Time = pc.ClockIn.Time.In(userLocation)
	}
	if pc.ClockOut.Valid {
		pc.ClockOut.Time = pc.ClockOut.Time.In(userLocation)
	}
	return
}

type PunchCardResponse struct {
	ClockIn   *time.Time `json:"clock_in"`
	ClockOut  *time.Time `json:"clock_out"`
	CreatedAt time.Time  `json:"created_at"`
}
