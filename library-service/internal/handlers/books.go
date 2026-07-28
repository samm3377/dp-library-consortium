package handlers

import (
	"encoding/json"
	"library-service/internal/models"
	"library-service/internal/service"
	"net/http"
	"strconv"
)

type BooksHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BooksHandler {
	return &BooksHandler{bookService}
}

func (b *BooksHandler) BooksHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var books []*models.Book
	var err error
	if title := r.FormValue("title"); title != "" {
		books, err = b.bookService.GetBooksByTitle(title)
	} else if author := r.FormValue("author"); author != "" {
		books, err = b.bookService.GetBooksByAuthor(author)
	} else if bookId := r.FormValue("bookId"); bookId != "" {
		id, err := strconv.Atoi(bookId)
		if err != nil {
			http.Error(w, "Parameter error", http.StatusBadRequest)
			return
		}
		books, err = b.bookService.GetBookByID(id)
	} else {
		books, err = b.bookService.GetAll()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ReservationRequest struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}

func (b *BooksHandler) ReservationHandler(w http.ResponseWriter, r *http.Request) {
	var res ReservationRequest
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, "Parameter error: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = b.bookService.ReserveBook(res.BookId)
	if err != nil {
		if err.Error() == "CHECK constraint failed: chk_books_available_copy" {
			http.Error(w, "No available books", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (b *BooksHandler) ReleaseHandler(w http.ResponseWriter, r *http.Request) {
	var res ReservationRequest
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, "Parameter error", http.StatusBadRequest)
		return
	}
	err = b.bookService.ReleaseBook(res.BookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
