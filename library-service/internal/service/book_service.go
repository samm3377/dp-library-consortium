package service

import (
	"library-service/internal/models"
	"library-service/internal/repository"
)

type BookService interface {
	GetBooksByTitle(title string) ([]*models.Book, error)
	GetBooksByAuthor(author string) ([]*models.Book, error)
	GetBookByID(id int) ([]*models.Book, error)
	GetAll() ([]*models.Book, error)
	ReserveBook(id int) error
	ReleaseBook(id int) error
}

type bookService struct {
	bookDB repository.BookRepository
}

func NewBookService(bookDB repository.BookRepository) BookService {
	return &bookService{bookDB}
}

func (b *bookService) GetBooksByAuthor(author string) ([]*models.Book, error) {
	return b.bookDB.FindByAuthor(author)
}

func (b *bookService) GetBookByID(id int) ([]*models.Book, error) {
	return b.bookDB.FindByID(id)
}

func (b *bookService) GetBooksByTitle(title string) ([]*models.Book, error) {
	return b.bookDB.FindByTitle(title)
}

func (b *bookService) GetAll() ([]*models.Book, error) {
	return b.bookDB.FindAll()
}

func (b *bookService) ReleaseBook(id int) error {
	return b.bookDB.IncreaseAvailability(id)
}

func (b *bookService) ReserveBook(id int) error {
	return b.bookDB.DecreaseAvailability(id)
}
