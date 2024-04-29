package handlers

import (
	"github.com/Alexande92/go-simple-library/internal/handlers/middleware"
	"net/http"
)

func RegisterRoutes(srv *http.ServeMux) {
	srv.HandleFunc("GET /api/v1/health", CheckHealth)

	srv.HandleFunc("GET /api/v1/books", GetBooks)
	srv.HandleFunc("GET /api/v1/book/{id}", GetBookById)
	srv.HandleFunc("DELETE /api/v1/book/{id}", DeleteBook)

	srv.Handle("POST /api/v1/book", middleware.SetDataValidation(SaveBook))
	srv.Handle("PUT /api/v1/book/{id}", middleware.SetDataValidation(UpdateBook))
}
