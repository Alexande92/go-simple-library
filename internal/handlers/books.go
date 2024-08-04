package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/entity"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"net/http"
	"strconv"
)

type BookHandler struct {
	db *storage.Storage
}

func NewBookHandler(db *storage.Storage) *BookHandler {
	return &BookHandler{
		db: db,
	}
}

func sendRequestError(w http.ResponseWriter, code int, response JsonErrorResponse) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
	}
}

func sendValidationError(w http.ResponseWriter, code int, response JsonValidationErrorResponse) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
	}
}

// TODO how to avoid adding JSON in each function?

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := h.db.GetAll()

	if err := json.NewEncoder(w).Encode(books); err != nil {
		sendRequestError(w, http.StatusInternalServerError, JsonErrorResponse{
			Message: "Internal server error " + err.Error(),
		})
		return
	}
}

func (h *BookHandler) SaveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book entity.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		sendRequestError(w, http.StatusBadRequest, JsonErrorResponse{
			Message: "Couldn't parse json",
		})
		return
	}

	validatedErrs := ValidateBook(book)

	if len(validatedErrs) != 0 {
		sendValidationError(w, http.StatusBadRequest, JsonValidationErrorResponse{
			Errors: validatedErrs,
		})

		return
	}

	book = h.db.Save(book)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		sendRequestError(w, http.StatusInternalServerError, JsonErrorResponse{
			Message: "Internal server error " + err.Error(),
		})
		return
	}
}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendRequestError(w, http.StatusBadRequest, JsonErrorResponse{
			Message: "Invalid book id",
		})
		return
	}

	book, err := h.db.GetById(int(bookId))

	if errors.Is(err, storage.ErrNotFound) {
		sendRequestError(w, http.StatusNotFound, JsonErrorResponse{
			Message: "Book not found",
		})
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendRequestError(w, http.StatusBadRequest, JsonErrorResponse{
			Message: "Invalid book id",
		})
		return
	}

	if err = h.db.Delete(int(bookId)); err != nil {
		sendRequestError(w, http.StatusInternalServerError, JsonErrorResponse{
			Message: "Internal server error " + err.Error(),
		})
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendRequestError(w, http.StatusBadRequest, JsonErrorResponse{
			Message: "Invalid book id",
		})
		return
	}

	var book entity.Book

	if err = json.NewDecoder(r.Body).Decode(&book); err != nil {
		sendRequestError(w, http.StatusBadRequest, JsonErrorResponse{
			Message: "Couldn't parse json",
		})
		return
	}

	validatedErrs := ValidateBook(book)

	if len(validatedErrs) != 0 {
		sendValidationError(w, http.StatusBadRequest, JsonValidationErrorResponse{
			Errors: validatedErrs,
		})

		return
	}
	book.ID = int(bookId)

	if err = h.db.Update(book); err != nil {
		sendRequestError(w, http.StatusNotFound, JsonErrorResponse{
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(book); err != nil {
		sendRequestError(w, http.StatusBadRequest, JsonErrorResponse{
			Message: "Internal error: " + err.Error(),
		})
	}
}
