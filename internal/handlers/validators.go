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

	err = validateLength(len(book.Author), -1, 255)
	errList = addValidationError(errList, "author", err)

	err = validateEmptiness(book.PublicationDate)
	errList = addValidationError(errList, "publicationDate", err)

	err = isEqual(len(book.PublicationDate), 7)
	errList = addValidationError(errList, "publicationDate", err)

	err = validateEmptiness(book.Title)
	errList = addValidationError(errList, "title", err)

	err = validateLength(len(book.Title), -1, 128)
	errList = addValidationError(errList, "title", err)

	err = validateLength(len(book.Publisher), 1, 255)
	errList = addValidationError(errList, "publisher", err)

	err = validateLength(len(book.Location), 1, 255)
	errList = addValidationError(errList, "location", err)

	err = validateLength(book.Edition, 1, -1)
	errList = addValidationError(errList, "edition", err)

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
	if len(val) == 0 {
		return errors.New("missing required field")
	}

	return nil
}

func validateLength(val int, min int, max int) error {
	if min != -1 {
		if val < min {
			return fmt.Errorf("field should contain more or equal than %d chars", min)
		}
	}

	if max != -1 {
		if val > max {
			return fmt.Errorf("field should contain less or equal than %d chars", max)
		}
	}

	return nil
}
