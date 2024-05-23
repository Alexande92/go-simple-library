package main

import (
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/handlers"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	fmt.Printf("Starting server at port %s\n", port)

	srv := http.NewServeMux()

	handlers.RegisterRoutes(srv)

	if err := http.ListenAndServe(port, srv); err != nil {
		log.Fatal(err)
	}
}
