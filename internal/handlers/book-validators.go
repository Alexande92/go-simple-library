package handlers

import (
	"errors"
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/storage"
)

type ValidationErrors struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Errors  []ErrorRes `json:"errors,omitempty"`
}

type ErrorRes struct {
	Field  string
	Reason string
}

func ValidateBookData(book storage.Book) []ErrorRes {
	errList := make([]ErrorRes, 0)

	err := validateEmptiness(book.Author)
	errList = addValidationError(errList, err)

	err = validateLength(len(book.Author), -1, 255)
	errList = addValidationError(errList, err)

	err = validateEmptiness(book.PublicationDate)
	errList = addValidationError(errList, err)

	err = isEqual(len(book.PublicationDate), 7)
	errList = addValidationError(errList, err)

	err = validateEmptiness(book.Title)
	errList = addValidationError(errList, err)

	err = validateLength(len(book.Title), -1, 128)
	errList = addValidationError(errList, err)

	err = validateLength(len(book.Publisher), 1, 255)
	errList = addValidationError(errList, err)

	err = validateLength(len(book.Location), 1, 255)
	errList = addValidationError(errList, err)

	err = validateLength(book.Edition, 1, -1)
	errList = addValidationError(errList, err)

	return errList
}

func addValidationError(errList []ErrorRes, err error) []ErrorRes {
	errList = append(errList, ErrorRes{
		Field:  "author",
		Reason: err.Error(),
	})

	return errList
}

func isEqual[T string | int](val T, toCompare T) error {
	if val != toCompare {
		err := fmt.Errorf("field should be equal to %s chars", toCompare)

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
			return fmt.Errorf("field should be contain more than %d chars", min)
		}
	}

	if min != -1 {
		if val > min {
			return fmt.Errorf("field should be contain less than %d chars", max)
		}
	}

	return nil
}
