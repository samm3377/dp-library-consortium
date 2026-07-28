package handlers

import (
	"central-server/internal/models"
	"central-server/internal/service"
	"central-server/internal/view"
	"log"
	"net/http"
	"strconv"
)

type Home struct {
	bookService        service.BookService
	render             view.Render
	reservationService service.ReservationService
}

func NewHomeHandler(bookService service.BookService, render view.Render, reservationService service.ReservationService) *Home {
	return &Home{bookService, render, reservationService}
}

type SearchPage struct {
	Query    string
	Books    []*models.Book
	Searched bool
	Username string
}

func (h *Home) Homehandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value("user").(models.User)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	findAll := r.URL.Query().Get("findAll")

	if query == "" {

		if findAll == "true" {

			books, errors := h.bookService.GetAll()
			if len(errors) != 0 {
				for _, err := range errors {
					log.Println(err.Error())
				}
			}

			data := SearchPage{
				Query:    "",
				Books:    books,
				Searched: true,
				Username: user.Username,
			}
			err := h.render.RenderTemplate(w, "search.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
		data := SearchPage{
			Query:    "",
			Books:    nil,
			Searched: false,
			Username: user.Username,
		}
		err := h.render.RenderTemplate(w, "search.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return

	} else {

		books, errors := h.bookService.FindBooks(query)
		if len(errors) != 0 {
			for _, err := range errors {
				log.Println(err.Error())
			}
		}

		data := SearchPage{
			Query:    query,
			Books:    books,
			Searched: true,
			Username: user.Username,
		}

		err := h.render.RenderTemplate(w, "search.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Home) ReservationHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid book id", http.StatusBadRequest)
		return
	}
	libraryID, err := strconv.Atoi(r.FormValue("libraryId"))
	if err != nil {
		http.Error(w, "Invalid library id", http.StatusBadRequest)
		return
	}

	h.reservationService.Reserve(bookID, libraryID, r.Context().Value("user").(models.User).ID)

	query := r.FormValue("query")
	if query != "" {
		http.Redirect(w, r, "/?query="+query, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/?findAll=true", http.StatusSeeOther)
	}
}
