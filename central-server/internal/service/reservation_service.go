package service

import (
	"bytes"
	"central-server/internal/models"
	"central-server/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ReservationService interface {
	Reserve(bookId int, libraryId int, userId int) error
	Release(bookId int, libraryId int, userId int) error
	ShowReservation(userID int) ([]*models.Book, error)
}

type reservationService struct {
	httpClient    *http.Client
	libraries     []*models.Library
	reservationDB repository.ReservationRepository
}

func NewReservationService(client http.Client, libraries []*models.Library, reservationDB repository.ReservationRepository) ReservationService {
	return &reservationService{&client, libraries, reservationDB}
}

type ReservationRequest struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}

// Reserve implements [ReservationService].
func (r *reservationService) Reserve(bookId int, libraryId int, userId int) error {
	baseURL := ""
	for _, l := range r.libraries {
		if l.ID == libraryId {
			baseURL = l.BaseURL
		}
	}

	req := ReservationRequest{
		BookId: bookId,
		UserId: userId,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := r.httpClient.Post(baseURL+"/reservation", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Reservation failed: %s", resp.Status)
	}

	reservation := models.Reservation{
		UserID:    userId,
		BookID:    bookId,
		LibraryID: libraryId,
	}

	err = r.reservationDB.AddReservation(reservation)

	return err
}

func (r *reservationService) Release(bookId int, libraryId int, userId int) error {
	baseURL := ""
	for _, l := range r.libraries {
		if l.ID == libraryId {
			baseURL = l.BaseURL
		}
	}

	req := ReservationRequest{
		BookId: bookId,
		UserId: userId,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := r.httpClient.Post(baseURL+"/release", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reservation failed: %s", resp.Status)
	}

	reservation := models.Reservation{
		UserID:    userId,
		BookID:    bookId,
		LibraryID: libraryId,
	}

	err = r.reservationDB.RemoveReservation(reservation)

	return err
}

func (r *reservationService) ShowReservation(userID int) ([]*models.Book, error) {
	reservations, err := r.reservationDB.GetReservationByUserID(userID)
	if err != nil {
		return nil, err
	}
	var books []*models.Book
	for _, res := range reservations {
		var booksResult []*models.Book
		baseURL := res.Library.BaseURL

		result, err := r.httpClient.Get(fmt.Sprintf("%s/books?bookId=%s", baseURL, url.QueryEscape(strconv.Itoa(res.BookID))))
		if err != nil {
			continue
		}
		defer result.Body.Close()
		err = json.NewDecoder(result.Body).Decode(&booksResult)
		if err != nil {
			return nil, err
		}
		booksResult[0].LibraryID = res.LibraryID
		booksResult[0].City = res.Library.City
		books = append(books, booksResult...)
	}
	return books, err
}
