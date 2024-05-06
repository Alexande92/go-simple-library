package handlers

import (
	"github.com/Alexande92/go-simple-library/internal/storage"
	"net/http"
)

func RegisterRoutes(srv *http.ServeMux) {
	db := storage.InitStorage()
	h := NewBaseHandler(db)

	srv.HandleFunc("GET /api/v1/health", CheckHealth)

	srv.HandleFunc("GET /api/v1/books", h.GetBooks)
	srv.HandleFunc("GET /api/v1/book/{id}", h.GetBookById)
	srv.HandleFunc("DELETE /api/v1/book/{id}", h.DeleteBook)

	srv.HandleFunc("POST /api/v1/book", h.SaveBook)
	srv.HandleFunc("PUT /api/v1/book/{id}", h.UpdateBook)
}
