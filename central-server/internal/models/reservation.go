package models

type Reservation struct {
	ID        int `gorm:"primaryKey"`
	UserID    int
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BookID    int
	LibraryID int
	Library   Library `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
