package models

type App struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"not null"`
	Secret string `gorm:"not null"`
}
