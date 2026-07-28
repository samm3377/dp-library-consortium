package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"central-server/internal/handlers"
	"central-server/internal/middleware"
	"central-server/internal/models"
	"central-server/internal/repository"
	"central-server/internal/service"
	"central-server/internal/view"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	config := LoadConfig()

	db, dbErr := gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{})
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Library{}, &models.Reservation{}); err != nil {
		log.Fatal(err)
	}

	var libraries []*models.Library
	if err := db.Find(&libraries).Error; err != nil {
		log.Fatal(err)
	}
	if len(libraries) == 0 {
		file, err := os.ReadFile("seed/libraries.json")
		if err != nil {
			log.Fatalf("Couldn't read libraries.json file: %v", err)
		}
		err = json.Unmarshal(file, &libraries)
		if err != nil {
			log.Fatalf("Parsing Error file libraries.json: %v", err)
		}
		err = db.Create(&libraries).Error
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	store := sessions.NewCookieStore([]byte(config.SecretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(7 * 24 * time.Hour.Seconds()),
	}

	userRepository := repository.NewUserRepository(db)
	reservationRepository := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepository)

	templates := LookupTemplates(config.TemplateDir)
	render := view.NewPageRender(templates)

	authHandler := handlers.NewAuthHandler(authService, store, render)

	r := mux.NewRouter()

	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir(config.TemplateDir))))

	r.HandleFunc("/register", authHandler.Register).Methods("Post", "Get")
	r.HandleFunc("/login", authHandler.Login).Methods("Post", "Get")
	r.HandleFunc("/logout", authHandler.Logout).Methods("Get")

	httpClient := http.Client{Timeout: 2 * time.Second}
	bookService := service.NewBookService(httpClient, libraries)
	reservationService := service.NewReservationService(httpClient, libraries, reservationRepository)
	searchHandler := handlers.NewHomeHandler(bookService, render, reservationService)
	profileHandler := handlers.NewProfileHandler(reservationService, render)

	r.HandleFunc("/", middleware.AuthMiddleware(store, searchHandler.Homehandler))
	r.HandleFunc("/reservation", middleware.AuthMiddleware(store, searchHandler.ReservationHandler))
	r.HandleFunc("/profile", middleware.AuthMiddleware(store, profileHandler.ProfileHandler))
	r.HandleFunc("/release", middleware.AuthMiddleware(store, profileHandler.ReleaseHandler))

	log.Printf("Server listening on port %s", config.ServerPort)

	server := &http.Server{
		Addr:           ":" + config.ServerPort,
		Handler:        r,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
