package models

type Book struct {
	ID            int
	Title         string
	Author        string
	AvailableCopy int
	LibraryID     int
	City          string
}
