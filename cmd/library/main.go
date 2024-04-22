package main

import (
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/http/handlers"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	fmt.Printf("Starting server at port %s\n", port)

	srv := &http.Server{
		Addr: port,
		// TODO: ask ho to move router to separate folder
		// TODO: and avoid pattern matching in main
		// TODO: Handler: router,
	}

	http.HandleFunc("/api/v1/health", handlers.CheckHealth)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
