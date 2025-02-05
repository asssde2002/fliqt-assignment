package models

type RoleName string

const (
	Admin RoleName = "admin"
	Staff RoleName = "staff"
)

func (r RoleName) Valid() bool {
	switch r {
	case Admin, Staff:
		return true
	default:
		return false
	}
}

type Role struct {
	ID   int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name RoleName `json:"name" gorm:"uniqueIndex;notNull;size:255"`
}
