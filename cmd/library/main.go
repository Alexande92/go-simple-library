package main

import (
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/handlers"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	fmt.Printf("Starting server at port %s\n", port)

	srv := http.NewServeMux()
	db := storage.NewStorage()

	handlers.RegisterRoutes(srv, db)

	if err := http.ListenAndServe(port, srv); err != nil {
		log.Fatal(err)
	}
}
