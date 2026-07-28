package repository

import (
	"central-server/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupReservationTestDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Library{}, &models.Reservation{})

	assert.NoError(t, err)
	return db
}

func TestReservationRepository_AddReservation_Success(t *testing.T) {

	db := setupReservationTestDB(t)
	repo := NewReservationRepository(db)
	reservation := models.Reservation{
		UserID:    1,
		BookID:    10,
		LibraryID: 5,
	}
	err := repo.AddReservation(reservation)

	assert.NoError(t, err)

	var saved models.Reservation
	err = db.First(&saved).Error

	assert.NoError(t, err)
	assert.Equal(t, 1, saved.UserID)
	assert.Equal(t, 10, saved.BookID)
	assert.Equal(t, 5, saved.LibraryID)
}

func TestReservationRepository_GetReservationByUserID_Success(t *testing.T) {

	db := setupReservationTestDB(t)
	repo := NewReservationRepository(db)
	library := models.Library{
		ID:      5,
		City:    "Rome",
		BaseURL: "http://library.com",
	}

	db.Create(&library)

	db.Create(&models.Reservation{
		UserID:    1,
		BookID:    10,
		LibraryID: 5,
	})

	db.Create(&models.Reservation{
		UserID:    2,
		BookID:    20,
		LibraryID: 5,
	})

	reservations, err := repo.GetReservationByUserID(1)

	assert.NoError(t, err)
	assert.Len(t, reservations, 1)
	assert.Equal(t, 1, reservations[0].UserID)
	assert.Equal(t, 10, reservations[0].BookID)
	assert.NotNil(t, reservations[0].Library)
	assert.Equal(t, "Rome", reservations[0].Library.City)
}

func TestReservationRepository_GetReservationByUserID_NoResults(t *testing.T) {

	db := setupReservationTestDB(t)
	repo := NewReservationRepository(db)
	reservations, err := repo.GetReservationByUserID(99)

	assert.NoError(t, err)
	assert.Empty(t, reservations)
}

func TestReservationRepository_RemoveReservation_Success(t *testing.T) {

	db := setupReservationTestDB(t)
	repo := NewReservationRepository(db)

	reservation := models.Reservation{
		UserID:    1,
		BookID:    10,
		LibraryID: 5,
	}

	db.Create(&reservation)
	err := repo.RemoveReservation(reservation)

	assert.NoError(t, err)

	var found models.Reservation
	result := db.First(&found)

	assert.Error(t, result.Error)
}

func TestReservationRepository_RemoveReservation_NotFound(t *testing.T) {

	db := setupReservationTestDB(t)
	repo := NewReservationRepository(db)

	reservation := models.Reservation{
		UserID:    1,
		BookID:    10,
		LibraryID: 5,
	}

	err := repo.RemoveReservation(reservation)

	assert.Error(t, err)
}
