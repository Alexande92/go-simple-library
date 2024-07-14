package handlers

import (
	"errors"
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/entities"
	"math"
)

type ValidationErrors struct {
	Errors []ErrorRes `json:"errors,omitempty"`
}

type ErrorRes struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func ValidateBook(book entities.Book) []ErrorRes {
	errList := make([]ErrorRes, 0)

	err := validateLength(book.Author, 1, math.MaxInt)
	errList = addValidationError(errList, "author", err)

	err = isEqual(len(book.PublicationDate), 7)
	errList = addValidationError(errList, "publicationDate", err)

	err = validateLength(book.Title, 1, 128)
	errList = addValidationError(errList, "title", err)

	err = validateLength(book.Publisher, 1, math.MaxInt)
	errList = addValidationError(errList, "publisher", err)

	err = validateLength(book.Location, 1, math.MaxInt)
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

func validateLength(val string, minLen, maxLen int) error {

	if len(val) < minLen {
		return fmt.Errorf("field should be at least %d characters long", minLen)
	}

	if len(val) > maxLen {
		return fmt.Errorf("field length should not exceed %d characters", maxLen)
	}

	return nil
}
