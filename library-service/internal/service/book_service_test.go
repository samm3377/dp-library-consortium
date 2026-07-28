package service

import (
	"errors"
	"library-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) FindByTitle(title string) ([]*models.Book, error) {
	args := m.Called(title)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByAuthor(author string) ([]*models.Book, error) {
	args := m.Called(author)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByID(id int) ([]*models.Book, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindAll() ([]*models.Book, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookRepository) IncreaseAvailability(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepository) DecreaseAvailability(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookService_GetBooksByTitle_Success(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)

	expected := []*models.Book{
		{
			Title: "The Hobbit",
		},
	}

	repo.On("FindByTitle", "Hobbit").Return(expected, nil)
	books, err := service.GetBooksByTitle("Hobbit")

	assert.NoError(t, err)
	assert.Equal(t, expected, books)
	repo.AssertExpectations(t)
}

func TestBookService_GetBooksByTitle_Error(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)
	expectedErr := errors.New("database error")
	repo.On("FindByTitle", "Hobbit").Return(nil, expectedErr)
	books, err := service.GetBooksByTitle("Hobbit")

	assert.Nil(t, books)
	assert.Equal(t, expectedErr, err)
}

func TestBookService_GetBooksByAuthor_Success(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)

	expected := []*models.Book{
		{
			Author: "Tolkien",
		},
	}

	repo.On("FindByAuthor", "Tolkien").Return(expected, nil)
	books, err := service.GetBooksByAuthor("Tolkien")

	assert.NoError(t, err)
	assert.Equal(t, expected, books)
}

func TestBookService_GetBookByID_Success(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)

	expected := []*models.Book{
		{
			ID: 10,
		},
	}

	repo.On("FindByID", 10).Return(expected, nil)
	books, err := service.GetBookByID(10)

	assert.NoError(t, err)
	assert.Equal(t, expected, books)
}

func TestBookService_GetAll_Success(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)

	expected := []*models.Book{
		{
			Title: "1984",
		},
	}

	repo.On("FindAll").Return(expected, nil)
	books, err := service.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expected, books)
}

func TestBookService_ReserveBook_Success(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)

	repo.On("DecreaseAvailability", 10).Return(nil)
	err := service.ReserveBook(10)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookService_ReserveBook_Error(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)
	expectedErr := errors.New("cannot reserve")
	repo.On("DecreaseAvailability", 10).Return(expectedErr)
	err := service.ReserveBook(10)

	assert.Equal(t, expectedErr, err)
}

func TestBookService_ReleaseBook_Success(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)
	repo.On("IncreaseAvailability", 10).Return(nil)
	err := service.ReleaseBook(10)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookService_ReleaseBook_Error(t *testing.T) {

	repo := new(MockBookRepository)
	service := NewBookService(repo)
	expectedErr := errors.New("cannot release")
	repo.On("IncreaseAvailability", 10).Return(expectedErr)
	err := service.ReleaseBook(10)

	assert.Equal(t, expectedErr, err)
}
