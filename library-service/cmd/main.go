package main

import (
	"encoding/json"
	"library-service/internal/handlers"
	"library-service/internal/models"
	"library-service/internal/repository"
	"library-service/internal/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, dbErr := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	if err := db.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal(err)
	}

	var books []*models.Book
	if err := db.Find(&books).Error; err != nil {
		log.Fatal(err)
	}
	if len(books) == 0 {
		booksFile := os.Getenv("BOOKS_FILE")
		file, err := os.ReadFile(booksFile)
		if err != nil {
			log.Fatalf("Couldn't read books.json file: %v", err)
		}
		err = json.Unmarshal(file, &books)
		if err != nil {
			log.Fatalf("Parsing Error file books.json: %v", err)
		}
		err = db.Create(&books).Error
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}

	bookRepository := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handlers.NewBookHandler(bookService)

	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.BooksHandler).Methods("GET")
	r.HandleFunc("/reservation", bookHandler.ReservationHandler).Methods("POST")
	r.HandleFunc("/release", bookHandler.ReleaseHandler).Methods("POST")

	log.Printf("Server listening on port %d", 8080)

	server := &http.Server{
		Addr:           ":" + "8080",
		Handler:        r,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(10 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
