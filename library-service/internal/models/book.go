package models

type Book struct {
	ID            int `gorm:"primaryKey"`
	Title         string
	Author        string
	AvailableCopy int `gorm:"check:available_copy >= 0"`
}
