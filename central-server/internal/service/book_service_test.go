package service

import (
	"central-server/internal/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/books", r.URL.Path)
		fmt.Fprintln(w, `[
			{
				"id":1,
				"title":"The Hobbit",
				"author":"Tolkien"
			}
		]`)
	}))
	defer server.Close()

	library := &models.Library{
		ID:      10,
		City:    "Rome",
		BaseURL: server.URL,
	}

	service := NewBookService(http.Client{}, []*models.Library{library})
	books, errs := service.GetAll()

	assert.Empty(t, errs)
	assert.Len(t, books, 1)
	assert.Equal(t, "The Hobbit", books[0].Title)
	assert.Equal(t, "Tolkien", books[0].Author)
	assert.Equal(t, 10, books[0].LibraryID)
	assert.Equal(t, "Rome", books[0].City)
}

func TestGetByTitle_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "The Hobbit", r.URL.Query().Get("title"))
		fmt.Fprintln(w, `[
			{
				"id":1,
				"title":"The Hobbit",
				"author":"Tolkien"
			}
		]`)
	}))
	defer server.Close()

	library := &models.Library{
		ID:      1,
		City:    "Milan",
		BaseURL: server.URL,
	}

	service := NewBookService(http.Client{}, []*models.Library{library})
	books, errs := service.GetByTitle("The Hobbit")

	assert.Empty(t, errs)
	assert.Len(t, books, 1)
	assert.Equal(t, "The Hobbit", books[0].Title)
	assert.Equal(t, 1, books[0].LibraryID)
	assert.Equal(t, "Milan", books[0].City)
}

func TestGetByAuthor_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "Tolkien", r.URL.Query().Get("author"))

		fmt.Fprintln(w, `[
			{
				"id":1,
				"title":"The Hobbit",
				"author":"Tolkien"
			}
		]`)
	}))
	defer server.Close()

	library := &models.Library{
		ID:      5,
		City:    "Turin",
		BaseURL: server.URL,
	}

	service := NewBookService(http.Client{}, []*models.Library{library})
	books, errs := service.GetByAuthor("Tolkien")

	assert.Empty(t, errs)
	assert.Len(t, books, 1)
	assert.Equal(t, "Tolkien", books[0].Author)
	assert.Equal(t, 5, books[0].LibraryID)
	assert.Equal(t, "Turin", books[0].City)
}

func TestFindBooks_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		author := r.URL.Query().Get("author")
		title := r.URL.Query().Get("title")

		if author == "Hobbit" {
			fmt.Fprintln(w, `[
				{
					"id":1,
					"title":"The Hobbit",
					"author":"Tolkien"
				}
			]`)
			return
		}

		if title == "Hobbit" {
			fmt.Fprintln(w, `[
				{
					"id":2,
					"title":"Hobbit Companion",
					"author":"Unknown"
				}
			]`)
			return
		}

		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	library := &models.Library{
		ID:      5,
		City:    "Turin",
		BaseURL: server.URL,
	}

	service := NewBookService(http.Client{}, []*models.Library{library})

	books, errs := service.FindBooks("Hobbit")

	assert.Empty(t, errs)
	assert.Len(t, books, 2)

	assert.Equal(t, "The Hobbit", books[0].Title)
	assert.Equal(t, "Tolkien", books[0].Author)
	assert.Equal(t, 5, books[0].LibraryID)
	assert.Equal(t, "Turin", books[0].City)

	assert.Equal(t, "Hobbit Companion", books[1].Title)
	assert.Equal(t, "Unknown", books[1].Author)
	assert.Equal(t, 5, books[1].LibraryID)
	assert.Equal(t, "Turin", books[1].City)
}

func TestFindBooks_Error(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	library := &models.Library{
		ID:      5,
		City:    "Turin",
		BaseURL: server.URL,
	}

	service := NewBookService(http.Client{}, []*models.Library{library})

	books, errs := service.FindBooks("Hobbit")

	assert.Empty(t, books)
	assert.NotEmpty(t, errs)
}

func TestGetAll_InvalidJSON(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `invalid json`)
	}))
	defer server.Close()

	library := &models.Library{
		ID:      1,
		City:    "Rome",
		BaseURL: server.URL,
	}

	service := NewBookService(http.Client{}, []*models.Library{library})
	books, errs := service.GetAll()

	assert.Empty(t, books)
	assert.Len(t, errs, 1)
}

func TestGetAll_HTTPError(t *testing.T) {

	library := &models.Library{
		ID:      1,
		City:    "Rome",
		BaseURL: "http://localhost:9999",
	}

	service := NewBookService(http.Client{}, []*models.Library{library})
	books, errs := service.GetAll()

	assert.Empty(t, books)
	assert.Len(t, errs, 1)
}

func TestAssignLibraryToBook(t *testing.T) {

	lib := &models.Library{
		ID:   12,
		City: "Turin",
	}

	books := []*models.Book{
		{
			Title: "Book1",
		},
		{
			Title: "Book2",
		},
	}

	assignLibraryToBook(lib, books)

	assert.Equal(t, 12, books[0].LibraryID)
	assert.Equal(t, "Turin", books[0].City)
	assert.Equal(t, 12, books[1].LibraryID)
	assert.Equal(t, "Turin", books[1].City)
}
