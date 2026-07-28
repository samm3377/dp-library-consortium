package service

import (
	"central-server/internal/models"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReservationRepository struct {
	mock.Mock
}

func (m *MockReservationRepository) AddReservation(r models.Reservation) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockReservationRepository) RemoveReservation(r models.Reservation) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockReservationRepository) GetReservationByUserID(userID int) ([]*models.Reservation, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Reservation), args.Error(1)
}

func TestReserve_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/reservation", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	mockRepo := new(MockReservationRepository)

	library := &models.Library{
		ID:      1,
		BaseURL: server.URL,
	}

	service := NewReservationService(http.Client{}, []*models.Library{library}, mockRepo)

	expected := models.Reservation{
		UserID:    3,
		BookID:    10,
		LibraryID: 1,
	}

	mockRepo.On("AddReservation", expected).Return(nil)
	err := service.Reserve(10, 1, 3)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReserve_ServerError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusBadRequest)
	}))
	defer server.Close()

	mockRepo := new(MockReservationRepository)

	library := &models.Library{
		ID:      1,
		BaseURL: server.URL,
	}

	service := NewReservationService(http.Client{}, []*models.Library{library}, mockRepo)
	err := service.Reserve(10, 1, 3)

	assert.Error(t, err)
}

func TestReserve_DBError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	mockRepo := new(MockReservationRepository)

	library := &models.Library{
		ID:      1,
		BaseURL: server.URL,
	}

	service := NewReservationService(http.Client{}, []*models.Library{library}, mockRepo)

	expected := models.Reservation{
		UserID:    3,
		BookID:    10,
		LibraryID: 1,
	}

	mockRepo.On("AddReservation", expected).Return(errors.New("db error"))
	err := service.Reserve(10, 1, 3)

	assert.Error(t, err)
}

func TestRelease_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/release", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mockRepo := new(MockReservationRepository)

	library := &models.Library{
		ID:      1,
		BaseURL: server.URL,
	}

	service := NewReservationService(http.Client{}, []*models.Library{library}, mockRepo)

	expected := models.Reservation{
		UserID:    3,
		BookID:    10,
		LibraryID: 1,
	}

	mockRepo.On("RemoveReservation", expected).Return(nil)
	err := service.Release(10, 1, 3)

	assert.NoError(t, err)
}

func TestShowReservation_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "10", r.URL.Query().Get("bookId"))

		fmt.Fprintln(w, `[
			{
				"id":10,
				"title":"The Hobbit",
				"author":"Tolkien"
			}
		]`)
	}))
	defer server.Close()

	mockRepo := new(MockReservationRepository)

	library := &models.Library{
		ID:      1,
		City:    "Rome",
		BaseURL: server.URL,
	}

	mockRepo.On("GetReservationByUserID", 5).Return([]*models.Reservation{
		{
			UserID:    5,
			BookID:    10,
			LibraryID: 1,
			Library:   *library,
		},
	}, nil)

	service := NewReservationService(http.Client{}, []*models.Library{library}, mockRepo)
	books, err := service.ShowReservation(5)

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "The Hobbit", books[0].Title)
	assert.Equal(t, 1, books[0].LibraryID)
	assert.Equal(t, "Rome", books[0].City)
}

func TestShowReservation_DBError(t *testing.T) {

	mockRepo := new(MockReservationRepository)
	mockRepo.On("GetReservationByUserID", 5).Return(nil, errors.New("db error"))
	service := NewReservationService(http.Client{}, nil, mockRepo)
	books, err := service.ShowReservation(5)

	assert.Nil(t, books)
	assert.Error(t, err)
}

func TestShowReservation_InvalidJSON(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{invalid json")
	}))
	defer server.Close()

	mockRepo := new(MockReservationRepository)

	library := models.Library{
		ID:      1,
		City:    "Rome",
		BaseURL: server.URL,
	}

	mockRepo.On("GetReservationByUserID", 5).Return([]*models.Reservation{
		{
			BookID:    10,
			LibraryID: 1,
			Library:   library,
		},
	}, nil)

	service := NewReservationService(http.Client{}, []*models.Library{&library}, mockRepo)
	books, err := service.ShowReservation(5)

	assert.Nil(t, books)
	assert.Error(t, err)
}
