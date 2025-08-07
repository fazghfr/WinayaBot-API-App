package model

type Task struct {
	ID     string `gorm:"type:char(36);primaryKey"`
	Title  string
	Status string
	UserID string `gorm:"type:char(36)"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Optional
}
