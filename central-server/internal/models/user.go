package models

import (
	"time"
)

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex"`
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
}
