package handlers

import (
	"errors"
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/storage"
)

type ValidationErrors struct {
	Errors []ErrorRes `json:"errors,omitempty"`
}

type ErrorRes struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func ValidateBook(book storage.Book) []ErrorRes {
	errList := make([]ErrorRes, 0)

	err := validateEmptiness(book.Author)
	errList = addValidationError(errList, "author", err)

	err = validateLength(book.Author, -1, 255)
	errList = addValidationError(errList, "author", err)

	err = validateEmptiness(book.PublicationDate)
	errList = addValidationError(errList, "publicationDate", err)

	err = isEqual(len(book.PublicationDate), 7)
	errList = addValidationError(errList, "publicationDate", err)

	err = validateEmptiness(book.Title)
	errList = addValidationError(errList, "title", err)

	err = validateLength(book.Title, -1, 128)
	errList = addValidationError(errList, "title", err)

	err = validateLength(book.Publisher, 1, 255)
	errList = addValidationError(errList, "publisher", err)

	err = validateLength(book.Location, 1, 255)
	errList = addValidationError(errList, "location", err)

	if book.Edition <= 0 {
		errList = addValidationError(errList, "edition", errors.New("edition should be a positive integer"))
	}

	return errList
}

func addValidationError(errList []ErrorRes, field string, err error) []ErrorRes {
	if err != nil {
		errList = append(errList, ErrorRes{
			Field:  field,
			Reason: err.Error(),
		})
	}

	return errList
}

func isEqual[T string | int](val T, toCompare T) error {
	if val != toCompare {
		err := fmt.Errorf("field should be equal to %v chars", toCompare)

		return err
	}

	return nil
}

func validateEmptiness(val string) error {
	if val == "" {
		return errors.New("missing required field")
	}

	return nil
}

func validateLength(val string, minLen, maxLen int) error {

	if len(val) < minLen {
		return fmt.Errorf("field should be at least %d characters long", minLen)
	}

	if len(val) > maxLen {
		return fmt.Errorf("field length should not exceed %d characters", maxLen)
	}

	return nil
}
