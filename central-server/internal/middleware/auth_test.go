package middleware

import (
	"central-server/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Success(t *testing.T) {

	store := sessions.NewCookieStore([]byte("secret-key"))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	session, _ := store.Get(req, "auth-session")
	session.Values["username"] = "test"
	session.Values["user_id"] = 10
	rec := httptest.NewRecorder()
	err := session.Save(req, rec)

	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Cookie", rec.Header().Get("Set-Cookie"))
	nextCalled := false

	handler := AuthMiddleware(store, func(w http.ResponseWriter, r *http.Request) {

		nextCalled = true
		user := r.Context().Value("user").(models.User)

		assert.Equal(t, "test", user.Username)
		assert.Equal(t, 10, user.ID)

		w.WriteHeader(http.StatusOK)
	})

	response := httptest.NewRecorder()
	handler(response, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestAuthMiddleware_MissingUsername(t *testing.T) {

	store := sessions.NewCookieStore([]byte("secret-key"))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	session, _ := store.Get(req, "auth-session")
	session.Values["user_id"] = 10
	rec := httptest.NewRecorder()
	session.Save(req, rec)
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Cookie", rec.Header().Get("Set-Cookie"))
	response := httptest.NewRecorder()
	handler := AuthMiddleware(store, func(w http.ResponseWriter, r *http.Request) {})
	handler(response, req)

	assert.Equal(t, http.StatusSeeOther, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func TestAuthMiddleware_MissingUserID(t *testing.T) {

	store := sessions.NewCookieStore([]byte("secret-key"))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	session, _ := store.Get(req, "auth-session")
	session.Values["username"] = "test"
	rec := httptest.NewRecorder()
	session.Save(req, rec)
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Cookie", rec.Header().Get("Set-Cookie"))
	response := httptest.NewRecorder()
	handler := AuthMiddleware(store, func(w http.ResponseWriter, r *http.Request) {})
	handler(response, req)

	assert.Equal(t, http.StatusSeeOther, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func TestAuthMiddleware_InvalidSession(t *testing.T) {

	store := sessions.NewCookieStore([]byte("secret-key"))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Cookie", "auth-session=invalid-cookie")
	response := httptest.NewRecorder()
	handler := AuthMiddleware(store, func(w http.ResponseWriter, r *http.Request) {})
	handler(response, req)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
