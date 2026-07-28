package handlers

import (
	"central-server/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProfileHandler_Unauthorized(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	handler.ProfileHandler(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProfileHandler_ShowReservationError(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	reservation.On("ShowReservation", 5).Return(nil, errors.New("database error"))
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ProfileHandler(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	reservation.AssertExpectations(t)
}

func TestProfileHandler_Success(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	books := []*models.Book{
		{
			Title: "Dune",
		},
	}

	reservation.On("ShowReservation", 5).Return(books, nil)
	render.On("RenderTemplate", mock.Anything, "profile.html", mock.AnythingOfType("handlers.SearchPage")).Return(nil)
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ProfileHandler(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	reservation.AssertExpectations(t)
	render.AssertExpectations(t)
}

func TestProfileHandler_RenderError(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	reservation.On("ShowReservation", 5).Return([]*models.Book{}, nil)
	render.On("RenderTemplate", mock.Anything, "profile.html", mock.Anything).Return(errors.New("template error"))

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ProfileHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestReleaseHandler_MethodNotAllowed(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	req := httptest.NewRequest(http.MethodGet, "/release", nil)
	rec := httptest.NewRecorder()
	handler.ReleaseHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestReleaseHandler_Unauthorized(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	req := httptest.NewRequest(http.MethodPost, "/release", nil)
	rec := httptest.NewRecorder()
	handler.ReleaseHandler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestReleaseHandler_InvalidBookID(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)
	handler := NewProfileHandler(reservation, render)
	form := url.Values{}
	form.Add("bookId", "abc")
	form.Add("libraryId", "1")
	req := httptest.NewRequest(http.MethodPost, "/release", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReleaseHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReleaseHandler_Success(t *testing.T) {

	reservation := new(MockReservationService)
	render := new(MockRender)

	handler := NewProfileHandler(reservation, render)

	reservation.On("Release", 10, 1, 5).Return(nil)
	form := url.Values{}
	form.Add("bookId", "10")
	form.Add("libraryId", "1")
	req := httptest.NewRequest(http.MethodPost, "/release", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = addUserToRequest(req)
	rec := httptest.NewRecorder()
	handler.ReleaseHandler(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/profile", rec.Header().Get("Location"))
	reservation.AssertExpectations(t)
}
