package handlers

import (
	"net/http"
	"time"

	"central-server/internal/models"
	"central-server/internal/service"
	"central-server/internal/view"

	"github.com/gorilla/sessions"
)

type Auth struct {
	authService service.AuthService
	store       *sessions.CookieStore
	render      view.Render
}

func NewAuthHandler(authService service.AuthService, store *sessions.CookieStore, render view.Render) *Auth {
	return &Auth{authService, store, render}
}

func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.render.RenderTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username:  r.FormValue("username"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		CreatedAt: time.Now(),
	}

	err := h.authService.Register(user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			http.Error(w, "Username already exists", http.StatusConflict)
		} else if err.Error() == "UNIQUE constraint failed: users.email" {
			http.Error(w, "Email already used", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.render.RenderTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var request struct {
		Username string
		Password string
	}

	request.Username = r.FormValue("username")
	request.Password = r.FormValue("password")

	user, err := h.authService.Login(request.Username, request.Password)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	session, _ := h.store.Get(r, "auth-session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "auth-session")
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
