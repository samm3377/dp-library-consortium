package middleware

import (
	"context"
	"log"
	"net/http"

	"central-server/internal/models"

	"github.com/gorilla/sessions"
)

func AuthMiddleware(store *sessions.CookieStore, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "auth-session")

		if err != nil {
			log.Printf("Error getting session: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["username"].(string)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user_id, ok := session.Values["user_id"].(int)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user := models.User{Username: username, ID: user_id}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
