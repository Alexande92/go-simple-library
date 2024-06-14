package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

// TODO how to avoid adding JSON in each function?

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := h.db.GetAll()

	//err := json.NewEncoder(w).Encode(books)

	if err := json.NewEncoder(w).Encode(books); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal error: " + err.Error())
		return
	}
}

func (h *BookHandler) SaveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book storage.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("Could not parse json")

		if err != nil {
			fmt.Println(err)
			//w.Write([]byte("{}"))
		}
		return
	}

	validatedErrs := ValidateBook(book)

	if len(validatedErrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationErrors{
			Errors: validatedErrs,
		})

		return
	}
	book = h.db.AddBook(book)
	book = h.db.Save(book)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid book id")
		return
	}

	book, err := h.db.GetById(int(bookId))

	if errors.Is(err, storage.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid book id")
		return
	}

	h.db.Delete(int(bookId))
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid book id")
		return
	}
	//_, err = h.db.GetById(int(bookId))
	//
	//if err != nil {
	//	w.WriteHeader(http.StatusNotFound)
	//	json.NewEncoder(w).Encode(err.Error())
	//	return
	//}

	var book storage.Book

	if err = json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not parse json")
		return
	}

	validatedErrs := ValidateBook(book)

	if len(validatedErrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationErrors{
			Errors: validatedErrs,
		})
		return
	}
	book.Id = int(bookId)

	err = h.db.Update(book)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(book)
	w.WriteHeader(http.StatusOK)
}
