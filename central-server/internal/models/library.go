package models

type Library struct {
	ID      int `gorm:"primaryKey"`
	City    string
	BaseURL string
}
