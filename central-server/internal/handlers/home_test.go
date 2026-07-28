package handlers

import (
	"bytes"
	"central-server/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHomeHandler_Unauthorized(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.Homehandler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHomeHandler_EmptySearch(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	render.On("RenderTemplate", mock.Anything, "search.html", mock.AnythingOfType("handlers.SearchPage")).Return(nil)
	handler.Homehandler(rec, req)

	render.AssertExpectations(t)
}

func TestHomeHandler_FindAll(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)

	books := []*models.Book{
		{
			Title: "Dune",
		},
	}

	bookService.On("GetAll").Return(books, []error(nil))
	req := httptest.NewRequest(http.MethodGet, "/?findAll=true", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	render.On("RenderTemplate", mock.Anything, "search.html", mock.AnythingOfType("handlers.SearchPage")).Return(nil)
	handler.Homehandler(rec, req)

	bookService.AssertExpectations(t)
	render.AssertExpectations(t)
}

func TestHomeHandler_SearchQuery(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)

	authorBooks := []*models.Book{
		{
			Title:  "The Hobbit",
			Author: "Tolkien",
		},
	}

	bookService.On("FindBooks", "Hobbit").Return(authorBooks, []error(nil))
	render.On("RenderTemplate", mock.Anything, "search.html", mock.AnythingOfType("handlers.SearchPage")).Return(nil)
	req := httptest.NewRequest(http.MethodGet, "/?query=Hobbit", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.Homehandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	bookService.AssertExpectations(t)
	render.AssertExpectations(t)
}

func TestHomeHandler_RenderError(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	render.On("RenderTemplate", mock.Anything, "search.html", mock.AnythingOfType("handlers.SearchPage")).Return(errors.New("template error"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.Homehandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestReservationHandler_MethodNotAllowed(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	req := httptest.NewRequest(http.MethodGet, "/reservation", nil)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestReservationHandler_Unauthorized(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	req := httptest.NewRequest(http.MethodPost, "/reservation", nil)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestReservationHandler_InvalidBookID(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBufferString("bookId=abc&libraryId=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReservationHandler_InvalidLibraryID(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBufferString("bookId=10&libraryId=abc"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReservationHandler_ServiceError(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	reservation.On("Reserve", 10, 1, 5).Return(errors.New("reservation error"))
	req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBufferString("bookId=10&libraryId=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	reservation.AssertExpectations(t)
}

func TestReservationHandler_RedirectWithQuery(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	reservation.On("Reserve", 10, 1, 5).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBufferString("bookId=10&libraryId=1&query=dune"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/?query=dune", rec.Header().Get("Location"))
}

func TestReservationHandler_RedirectFindAll(t *testing.T) {

	bookService := new(MockBookService)
	render := new(MockRender)
	reservation := new(MockReservationService)
	handler := NewHomeHandler(bookService, render, reservation)
	reservation.On("Reserve", 10, 1, 5).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBufferString("bookId=10&libraryId=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReservationHandler(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/?findAll=true", rec.Header().Get("Location"))
}
