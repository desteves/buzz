package main

import (
	"log"
	"net/http"
)

// main is the entry point for the application.
func main() {
	// Create a new server and set the handler.
	server := &http.Server{
		Addr:    ":8000",
		Handler: New(),
	}

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}
