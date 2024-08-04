package handlers

import (
	"github.com/Alexande92/go-simple-library/internal/storage"
	"net/http"
)

func RegisterRoutes(srv *http.ServeMux, db *storage.Storage) {
	h := NewBookHandler(db)

	srv.HandleFunc("GET /api/v1/health", CheckHealth)

	srv.HandleFunc("GET /api/v1/books", h.GetBooks)
	srv.HandleFunc("GET /api/v1/books/{id}", h.GetBookById)
	srv.HandleFunc("DELETE /api/v1/books/{id}", h.DeleteBook)

	srv.HandleFunc("POST /api/v1/books", h.SaveBook)
	srv.HandleFunc("PUT /api/v1/books/{id}", h.UpdateBook)
}
