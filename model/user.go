package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}
