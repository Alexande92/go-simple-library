package handlers

import (
	"encoding/json"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"io"
	"net/http"
	"strconv"
)

// TODO how to avoid adding JSON in each function?

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := storage.GetAll()

	err := json.NewEncoder(w).Encode(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal error")
		return
	}
}

func SaveBook(w http.ResponseWriter, r *http.Request, book storage.Book) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not parse json")
		return
	}

	book.Save()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not fetch book, wrong id format")
		return
	}

	book, ok := storage.GetById(bookId)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Book not found")
		return
	}

	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not fetch book")
		return
	}
	storage.Delete(bookId)
}

func UpdateBook(w http.ResponseWriter, r *http.Request, book storage.Book) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not fetch book")
		return
	}

	_, ok := storage.GetById(bookId)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Book not found")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&book)

	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not parse json")
		return
	}

	book.Id = bookId

	book.Update()

	w.WriteHeader(http.StatusOK)
}
