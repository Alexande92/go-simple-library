package handlers

import (
	"github.com/Alexande92/go-simple-library/internal/storage"
)

type ValidationErrorRes struct {
	Code    int               `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func ValidateBookData(book storage.Book) map[string]string {
	errList := make(map[string]string)

	if book.Author == "" {
		errList["Author"] = "missing required field"
	} else if len(book.Author) > 255 {
		errList["Author"] = "field more than 255 chars"
	}

	if book.PublicationDate == "" {
		errList["PublicationDate"] = "missing required field"
	} else if len(book.PublicationDate) != 7 {
		errList["PublicationDate"] = "field Date should contain 7 chars"
	}

	if book.Title == "" {
		errList["Title"] = "missing required field"
	} else if len(book.Title) > 128 {
		errList["Title"] = "field more than 128 chars"
	}

	if len(book.Publisher) < 1 || len(book.Publisher) > 255 {
		errList["Publisher"] = "field should contains from 1 to 255 chars"
	}

	if len(book.Location) < 1 || len(book.Location) > 255 {
		errList["Publisher"] = "field should contains from 1 to 255 chars"
	}

	if book.Edition < 1 {
		errList["Publisher"] = "field should be at least equal to 1"
	}

	return errList
}
