package service

import (
	"central-server/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type BookService interface {
	GetAll() ([]*models.Book, []error)
	GetByTitle(title string) ([]*models.Book, []error)
	GetByAuthor(author string) ([]*models.Book, []error)
	FindBooks(input string) ([]*models.Book, []error)
}

type bookService struct {
	httpClient *http.Client
	libraries  []*models.Library
}

func NewBookService(client http.Client, libraries []*models.Library) BookService {
	return &bookService{&client, libraries}
}

func (b *bookService) GetAll() ([]*models.Book, []error) {
	result := []*models.Book{}
	var errors []error
	for _, lib := range b.libraries {
		books := []*models.Book{}
		resp, err := b.httpClient.Get(fmt.Sprintf("%s/books", lib.BaseURL))
		if err != nil {
			errors = append(errors, err)
			continue
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&books)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		assignLibraryToBook(lib, books)
		result = append(result, books...)
	}
	return result, errors
}

func (b *bookService) GetByAuthor(author string) ([]*models.Book, []error) {
	result := []*models.Book{}
	var errors []error
	for _, lib := range b.libraries {
		books := []*models.Book{}
		resp, err := b.httpClient.Get(fmt.Sprintf("%s/books?author=%s", lib.BaseURL, url.QueryEscape(author)))
		if err != nil {
			errors = append(errors, err)
			continue
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&books)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		assignLibraryToBook(lib, books)
		result = append(result, books...)
	}
	return result, errors
}

func (b *bookService) GetByTitle(title string) ([]*models.Book, []error) {
	result := []*models.Book{}
	var errors []error
	for _, lib := range b.libraries {
		books := []*models.Book{}
		resp, err := b.httpClient.Get(fmt.Sprintf("%s/books?title=%s", lib.BaseURL, url.QueryEscape(title)))
		if err != nil {
			errors = append(errors, err)
			continue
		}
		err = json.NewDecoder(resp.Body).Decode(&books)
		resp.Body.Close()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		assignLibraryToBook(lib, books)
		result = append(result, books...)
	}
	return result, errors
}

func (b *bookService) FindBooks(input string) ([]*models.Book, []error) {
	booksByAuthor, errorsByAuthor := b.GetByAuthor(input)
	booksByTitle, errorsByTitle := b.GetByTitle(input)

	books := append(booksByAuthor, booksByTitle...)
	errors := append(errorsByAuthor, errorsByTitle...)
	books = removeDuplicates(books)
	return books, errors
}

func removeDuplicates(books []*models.Book) []*models.Book {
	type key struct {
		ID      int
		Library int
	}

	seen := make(map[key]bool)
	result := make([]*models.Book, 0)

	for _, book := range books {
		k := key{
			ID:      book.ID,
			Library: book.LibraryID,
		}

		if !seen[k] {
			seen[k] = true
			result = append(result, book)
		}
	}
	return result
}

func assignLibraryToBook(library *models.Library, books []*models.Book) {
	for _, book := range books {
		book.LibraryID = library.ID
		book.City = library.City
	}
}
