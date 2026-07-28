package handlers

import (
	"bytes"
	"errors"
	"library-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) GetBooksByTitle(title string) ([]*models.Book, error) {
	args := m.Called(title)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookService) GetBooksByAuthor(author string) ([]*models.Book, error) {
	args := m.Called(author)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookService) GetBookByID(id int) ([]*models.Book, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookService) GetAll() ([]*models.Book, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookService) ReserveBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookService) ReleaseBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBooksHandler_GetAll(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)

	books := []*models.Book{
		{
			Title: "Dune",
		},
	}

	service.On("GetAll").Return(books, nil)
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	handler.BooksHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	service.AssertExpectations(t)
}

func TestBooksHandler_ByTitle(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("GetBooksByTitle", "Dune").Return([]*models.Book{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/books?title=Dune", nil)
	rec := httptest.NewRecorder()
	handler.BooksHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBooksHandler_ByAuthor(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("GetBooksByAuthor", "Tolkien").Return([]*models.Book{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/books?author=Tolkien", nil)
	rec := httptest.NewRecorder()
	handler.BooksHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBooksHandler_ByID(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("GetBookByID", 10).Return([]*models.Book{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/books?bookId=10", nil)
	rec := httptest.NewRecorder()
	handler.BooksHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBooksHandler_InvalidID(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/books?bookId=abc", nil)
	rec := httptest.NewRecorder()
	handler.BooksHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBooksHandler_ServiceError(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("GetAll").Return(nil, errors.New("database error"))
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	handler.BooksHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestReservationHandler_Success(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("ReserveBook", 10).Return(nil)

	body := bytes.NewBufferString(`
		{
			"bookId":10,
			"userId":5
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/reservation", body)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestReservationHandler_InvalidJSON(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBufferString("invalid"))
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReservationHandler_NoBooks(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("ReserveBook", 10).Return(errors.New("CHECK constraint failed: chk_books_available_copy"))

	body := bytes.NewBufferString(`
		{
			"bookId":10
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/reservation", body)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestReleaseHandler_Success(t *testing.T) {

	service := new(MockBookService)
	handler := NewBookHandler(service)
	service.On("ReleaseBook", 10).Return(nil)

	body := bytes.NewBufferString(`
		{
			"bookId":10
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/release", body)
	rec := httptest.NewRecorder()
	handler.ReleaseHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
