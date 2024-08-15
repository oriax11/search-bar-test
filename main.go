package main

import (
	"fmt"
	"log"
	"net/http"
	groupie "Groupie/handlers"
)

func main() {
	// Set up route handlers
	http.HandleFunc("/", groupie.ArtistsHandler)
	http.HandleFunc("/artist/", groupie.ArtistDetailHandler)

	// Print server information
	fmt.Printf("Server listening on port %s...\n", groupie.Port)
	fmt.Printf("Visit http://localhost:%s to view the application\n", groupie.Port)

	// Start the HTTP server
	log.Printf("Starting server on :%s ", groupie.Port)
	err := http.ListenAndServe(":"+groupie.Port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
