package model

type User struct {
	ID         string `gorm:"type:char(36);primaryKey"`
	Discord_id string `gorm:"unique"`
	Tasks      []Task `gorm:"foreignKey:UserID"` // One-to-Many relationship
}
