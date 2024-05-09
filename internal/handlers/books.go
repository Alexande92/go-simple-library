package handlers

import (
	"encoding/json"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"io"
	"net/http"
	"strconv"
)

type BookHandler struct {
	db storage.Storage
}

func NewBookHandler(db storage.Storage) *BookHandler {
	return &BookHandler{
		db: db,
	}
}

// TODO how to avoid adding JSON in each function?

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	table := storage.FromTable(h.db, "book")
	books := table.GetAll()

	err := json.NewEncoder(w).Encode(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal error: " + err.Error())
		return
	}

	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) SaveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book storage.Book
	table := storage.FromTable(h.db, "book")

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not parse json")
		return
	}

	validatedErrs := ValidateBookData(book)

	if len(validatedErrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ValidationErrors{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  validatedErrs,
		})

		return
	}

	table.Save(&book)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid book id")
		return
	}

	table := storage.FromTable(h.db, "book")
	book, err := table.GetById(int(bookId))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid book id")
		return
	}

	table := storage.FromTable(h.db, "book")
	table.Delete(int(bookId))
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid book id")
		return
	}
	table := storage.FromTable(h.db, "book")

	_, err = table.GetById(int(bookId))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	var book storage.Book

	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not parse json")
		return
	}

	validatedErrs := ValidateBookData(book)

	if len(validatedErrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationErrors{
			Code:    400,
			Message: "Validation failed",
			Errors:  validatedErrs,
		})
		return
	}
	book.Id = int(bookId)

	table.Update(&book)

	w.WriteHeader(http.StatusOK)
}
