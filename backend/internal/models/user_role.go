package models

type UserRole struct {
	UserID int64 `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	RoleID int64 `json:"role_id" gorm:"primaryKey;autoIncrement:false"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}
