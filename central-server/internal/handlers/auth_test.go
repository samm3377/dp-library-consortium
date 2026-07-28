package handlers

import (
	"bytes"
	"central-server/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_Register_GET(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	render.On("RenderTemplate", mock.Anything, "register.html", nil).Return(nil)
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	render.AssertExpectations(t)
}

func TestAuth_Register_Success(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	auth.On("Register", "test", "test@test.com", "password", mock.Anything).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("username=test&email=test@test.com&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/login", rec.Header().Get("Location"))
}

func TestAuth_Register_DuplicateUsername(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	auth.On("Register", "test", "mail@test.com", "123", mock.Anything).Return(errors.New("UNIQUE constraint failed: users.username"))
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("username=test&email=mail@test.com&password=123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAuth_Register_DuplicateEmail(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	auth.On("Register", "test", "mail@test.com", "123", mock.Anything).Return(errors.New("UNIQUE constraint failed: users.email"))
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("username=test&email=mail@test.com&password=123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAuth_Login_GET(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	render.On("RenderTemplate", mock.Anything, "login.html", nil).Return(nil)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuth_Login_Success(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	auth.On("Login", "test", "123").Return(&models.User{ID: 1, Username: "test"}, nil)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("username=test&password=123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.Login(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	assert.NotEmpty(t, rec.Header().Get("Set-Cookie"))
}

func TestAuth_Login_InvalidCredentials(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	auth.On("Login", "test", "123").Return(nil, errors.New("record not found"))
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("username=test&password=123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.Login(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuth_Logout(t *testing.T) {

	auth := new(MockAuthService)
	render := new(MockRender)
	store := sessions.NewCookieStore([]byte("secret"))
	handler := NewAuthHandler(auth, store, render)
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rec := httptest.NewRecorder()
	handler.Logout(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/login", rec.Header().Get("Location"))
}
