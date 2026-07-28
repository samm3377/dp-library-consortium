package repository

import (
	"library-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBookTestDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Book{})

	assert.NoError(t, err)

	return db
}

func TestBookRepository_FindByTitle_Success(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)

	db.Create(&models.Book{
		Title:         "The Hobbit",
		Author:        "Tolkien",
		AvailableCopy: 3,
	})

	db.Create(&models.Book{
		Title:         "Harry Potter",
		Author:        "Rowling",
		AvailableCopy: 5,
	})

	books, err := repo.FindByTitle("Hobbit")

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "The Hobbit", books[0].Title)
}

func TestBookRepository_FindByAuthor_Success(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)

	db.Create(&models.Book{
		Title:  "Book 1",
		Author: "Tolkien",
	})

	db.Create(&models.Book{
		Title:  "Book 2",
		Author: "Rowling",
	})

	books, err := repo.FindByAuthor("Tolkien")

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Tolkien", books[0].Author)
}

func TestBookRepository_FindByID_Success(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)

	book := models.Book{
		Title:  "Dune",
		Author: "Frank Herbert",
	}

	db.Create(&book)
	books, err := repo.FindByID(int(book.ID))

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Dune", books[0].Title)
}

func TestBookRepository_FindAll_Success(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)

	db.Create(&models.Book{
		Title: "Book 1",
	})

	db.Create(&models.Book{
		Title: "Book 2",
	})

	books, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, books, 2)
}

func TestBookRepository_FindAll_Empty(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)
	books, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Empty(t, books)
}

func TestBookRepository_IncreaseAvailability_Success(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)

	book := models.Book{
		Title:         "1984",
		AvailableCopy: 2,
	}

	db.Create(&book)
	err := repo.IncreaseAvailability(int(book.ID))

	assert.NoError(t, err)

	var updated models.Book
	db.First(&updated, book.ID)

	assert.Equal(t, 3, updated.AvailableCopy)
}

func TestBookRepository_DecreaseAvailability_Success(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)

	book := models.Book{
		Title:         "1984",
		AvailableCopy: 2,
	}

	db.Create(&book)
	err := repo.DecreaseAvailability(int(book.ID))

	assert.NoError(t, err)

	var updated models.Book
	db.First(&updated, book.ID)

	assert.Equal(t, 1, updated.AvailableCopy)
}

func TestBookRepository_IncreaseAvailability_NotFound(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)
	err := repo.IncreaseAvailability(999)

	assert.Error(t, err)
}

func TestBookRepository_DecreaseAvailability_NotFound(t *testing.T) {

	db := setupBookTestDB(t)
	repo := NewBookRepository(db)
	err := repo.DecreaseAvailability(999)

	assert.Error(t, err)
}
