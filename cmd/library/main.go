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

	srv := &http.Server{
		Addr: port,
	}

	http.HandleFunc("/api/v1/health", handlers.CheckHealth)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
