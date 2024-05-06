package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"io"
	"net/http"
	"strconv"
)

type BaseHandler struct {
	db *storage.Storage
}

func NewBaseHandler(db *storage.Storage) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

// TODO how to avoid adding JSON in each function?

func (h *BaseHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	table := storage.FromTable(*h.db, "book")
	books := table.GetAll()

	err := json.NewEncoder(w).Encode(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal error")
		return
	}
}

func (h *BaseHandler) SaveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book storage.Book
	table := storage.FromTable(*h.db, "book")

	err := json.NewDecoder(r.Body).Decode(&book)
	//err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not parse json")
		return
	}

	validatedErrs := ValidateBookData(book)
	fmt.Println(validatedErrs)

	if len(validatedErrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ValidationErrorRes{
			Code:    400,
			Message: "Validation failed",
			Errors:  validatedErrs,
		})

		fmt.Println(err)
		return
	}

	table.Save(&book)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BaseHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not fetch book, wrong id format")
		return
	}

	table := storage.FromTable(*h.db, "book")
	book, ok := table.GetById(bookId)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Book not found")
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *BaseHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not fetch book")
		return
	}

	table := storage.FromTable(*h.db, "book")
	table.Delete(bookId)
}

func (h *BaseHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not fetch book")
		return
	}
	table := storage.FromTable(*h.db, "book")

	_, ok := table.GetById(bookId)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Book not found")
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
		json.NewEncoder(w).Encode(ValidationErrorRes{
			Code:    400,
			Message: "Validation failed",
			Errors:  validatedErrs,
		})
		return
	}
	book.Id = bookId

	table.Update(&book)

	w.WriteHeader(http.StatusOK)
}
