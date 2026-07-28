package repository

import (
	"central-server/internal/models"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	AddReservation(reservation models.Reservation) error
	RemoveReservation(reservation models.Reservation) error
	GetReservationByUserID(userID int) ([]*models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

func (r *reservationRepository) AddReservation(reservation models.Reservation) error {
	return r.db.Create(&reservation).Error
}

func (r *reservationRepository) GetReservationByUserID(userID int) ([]*models.Reservation, error) {
	var reservations []*models.Reservation
	err := r.db.Preload("Library").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) RemoveReservation(reservation models.Reservation) error {
	var found models.Reservation

	err := r.db.
		Where("user_id = ? AND book_id = ? AND library_id = ?",
			reservation.UserID,
			reservation.BookID,
			reservation.LibraryID,
		).
		First(&found).Error

	if err != nil {
		return err
	}
	return r.db.Delete(&found).Error
}
