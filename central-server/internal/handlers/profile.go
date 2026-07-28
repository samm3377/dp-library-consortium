package handlers

import (
	"central-server/internal/models"
	"central-server/internal/service"
	"central-server/internal/view"
	"net/http"
	"strconv"
)

type Profile struct {
	reservationService service.ReservationService
	render             view.Render
}

func NewProfileHandler(reservationService service.ReservationService, render view.Render) *Profile {
	return &Profile{reservationService, render}
}

func (p *Profile) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	books, err := p.reservationService.ShowReservation(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := SearchPage{
		Books:    books,
		Username: user.Username,
	}

	err = p.render.RenderTemplate(w, "profile.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Profile) ReleaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, ok := r.Context().Value("user").(models.User)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookID, err := strconv.Atoi(r.FormValue("bookId"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	libraryID, err := strconv.Atoi(r.FormValue("libraryId"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = p.reservationService.Release(bookID, libraryID, r.Context().Value("user").(models.User).ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
