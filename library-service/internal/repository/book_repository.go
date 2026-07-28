package repository

import (
	"library-service/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	FindByTitle(title string) ([]*models.Book, error)
	FindByAuthor(author string) ([]*models.Book, error)
	FindByID(id int) ([]*models.Book, error)
	FindAll() ([]*models.Book, error)
	IncreaseAvailability(id int) error
	DecreaseAvailability(id int) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (b *bookRepository) FindByAuthor(author string) ([]*models.Book, error) {
	var books []*models.Book
	err := b.db.Where("author like ?", "%"+author+"%").Find(&books).Error
	return books, err
}

func (b *bookRepository) FindByTitle(title string) ([]*models.Book, error) {
	var books []*models.Book
	err := b.db.Where("title like ?", "%"+title+"%").Find(&books).Error
	return books, err
}

func (b *bookRepository) FindByID(id int) ([]*models.Book, error) {
	var books []*models.Book
	err := b.db.Where("id = ?", id).Find(&books).Error
	return books, err
}

func (b *bookRepository) FindAll() ([]*models.Book, error) {
	var books []*models.Book
	err := b.db.Find(&books).Error
	return books, err
}

func (b *bookRepository) IncreaseAvailability(id int) error {
	var book models.Book
	err := b.db.Where("id = ?", id).First(&book).Error
	if err != nil {
		return err
	}
	book.AvailableCopy = book.AvailableCopy + 1
	return b.db.Model(&book).Update("available_copy", book.AvailableCopy).Error
}

func (b *bookRepository) DecreaseAvailability(id int) error {
	var book models.Book
	err := b.db.Where("id = ?", id).First(&book).Error
	if err != nil {
		return err
	}
	book.AvailableCopy = book.AvailableCopy - 1
	return b.db.Model(&book).Update("available_copy", book.AvailableCopy).Error
}
