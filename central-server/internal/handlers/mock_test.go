package handlers

import (
	"central-server/internal/models"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(
	username string,
	email string,
	password string,
	createdAt time.Time,
) error {

	args := m.Called(username, email, password, mock.Anything)

	return args.Error(0)
}

func (m *MockAuthService) Login(
	username string,
	password string,
) (*models.User, error) {

	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

type MockBookService struct {
	mock.Mock
}

// FindBooks implements [service.BookService].
func (m *MockBookService) FindBooks(input string) ([]*models.Book, []error) {
	args := m.Called(input)
	var books []*models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]*models.Book)
	}

	var errs []error
	if args.Get(1) != nil {
		errs = args.Get(1).([]error)
	}
	return books, errs
}

func (m *MockBookService) GetAll() ([]*models.Book, []error) {

	args := m.Called()
	var books []*models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]*models.Book)
	}

	var errs []error
	if args.Get(1) != nil {
		errs = args.Get(1).([]error)
	}
	return books, errs
}

func (m *MockBookService) GetByTitle(title string) ([]*models.Book, []error) {

	args := m.Called(title)
	var books []*models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]*models.Book)
	}

	var errs []error
	if args.Get(1) != nil {
		errs = args.Get(1).([]error)
	}
	return books, errs
}

func (m *MockBookService) GetByAuthor(author string) ([]*models.Book, []error) {

	args := m.Called(author)
	var books []*models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]*models.Book)
	}

	var errs []error
	if args.Get(1) != nil {
		errs = args.Get(1).([]error)
	}
	return books, errs
}

type MockReservationService struct {
	mock.Mock
}

func (m *MockReservationService) Reserve(
	bookID int,
	libraryID int,
	userID int,
) error {

	args := m.Called(bookID, libraryID, userID)
	return args.Error(0)
}

func (m *MockReservationService) Release(
	bookID int,
	libraryID int,
	userID int,
) error {

	args := m.Called(bookID, libraryID, userID)
	return args.Error(0)
}

func (m *MockReservationService) ShowReservation(
	userID int,
) ([]*models.Book, error) {

	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Book), args.Error(1)
}

type MockRender struct {
	mock.Mock
}

func (m *MockRender) RenderTemplate(
	w io.Writer,
	name string,
	data any,
) error {

	args := m.Called(w, name, data)
	return args.Error(0)
}

func addUserToRequest(req *http.Request) *http.Request {

	user := models.User{
		ID:       5,
		Username: "test",
	}

	ctx := context.WithValue(req.Context(), "user", user)
	return req.WithContext(ctx)
}
