package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	PassHash []uint8
	Is_Admin bool `gorm:"not null"`
}
