package middleware

import (
	"encoding/json"
	"errors"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"net/http"
)

type ValidatedHandler func(http.ResponseWriter, *http.Request, storage.Book)
type ValidateData struct {
	handler ValidatedHandler
}

// TODO make this middleware not only Book related, more universal

func (ea *ValidateData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var errList error
	var book storage.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Wrong format of json data")
		return
	}

	if book.Author == "" {
		errList = errors.Join(errList, errors.New("missing required field Author"))
	} else if len(book.Author) > 255 {
		errList = errors.Join(errList, errors.New("field Author more than 255 chars"))
	}

	if book.PublicationDate == "" {
		errList = errors.Join(errList, errors.New("missing required field Publication Date"))
	} else if len(book.PublicationDate) != 7 {
		errList = errors.Join(errList, errors.New("field Publication Date should contain 7 chars"))
	}

	if book.Title == "" {
		errList = errors.Join(errList, errors.New("missing required field Title"))
	} else if len(book.Title) > 128 {
		errList = errors.Join(errList, errors.New("field Title more than 128 chars"))
	}

	if len(book.Publisher) < 1 || len(book.Publisher) > 255 {
		errList = errors.Join(errList, errors.New("field Publisher should contains from 1 to 255 chars"))
	}

	if len(book.Location) < 1 || len(book.Location) > 255 {
		errList = errors.Join(errList, errors.New("field Location should contains from 1 to 255 chars"))
	}

	if book.Edition < 1 {
		errList = errors.Join(errList, errors.New("field Edition should be at least equal to 1"))
	}

	if errList != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errList.Error())
		return
	}

	ea.handler(w, r, book)
}

func SetDataValidation(handlerToWrap ValidatedHandler) *ValidateData {
	return &ValidateData{handlerToWrap}
}
