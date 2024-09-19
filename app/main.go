package main

import (
	"log"
	"net/http"
	"os"

	"main/web"
)

// main is the entry point for the application.
func main() {
	h := web.NewHandler()
	// Create a new server and set the handler.
	server := &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"), // "localhost:8000"
		Handler: h.Handler,
	}
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	}
}
