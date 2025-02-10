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
	CreatedAt time.Time    `json:"created_at" gorm:"notNull"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (pc *PunchCard) BeforeCreate(db *gorm.DB) (err error) {
	currTime := time.Now().UTC()
	if pc.CreatedAt.IsZero() {
		pc.CreatedAt = currTime
	}
	if !pc.ClockIn.Valid {
		pc.ClockIn = sql.NullTime{Time: currTime, Valid: true}
	}
	return
}

func (pc *PunchCard) AfterFind(db *gorm.DB) (err error) {
	pc.CreatedAt = pc.CreatedAt.Local()
	if pc.ClockIn.Valid {
		pc.ClockIn.Time = pc.ClockIn.Time.Local()
	}
	if pc.ClockOut.Valid {
		pc.ClockOut.Time = pc.ClockOut.Time.Local()
	}
	return
}

type PunchCardResponse struct {
	ClockIn   *time.Time `json:"clock_in"`
	ClockOut  *time.Time `json:"clock_out"`
	CreatedAt time.Time  `json:"created_at"`
}
