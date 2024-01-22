// Package declaration, main package is the entry point of the Go program
package main

// Importing necessary packages
import (
	"encoding/json" // Package for encoding and decoding JSON data
	"log"           // Package for logging
	"net/http"      // Package for handling HTTP requests

	"github.com/gorilla/mux" // External package for routing HTTP requests
)

// Main function, the entry point of the program
func main() {
	// Creating a new router using the gorilla/mux package
	router := mux.NewRouter()

	// Handling requests to the root path ("/") with a function
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// Setting the Content-Type header to indicate that the response is JSON
		rw.Header().Set("Content-Type", "application/json")

		// Encoding a JSON response with a map containing a "data" key and a message value
		json.NewEncoder(rw).Encode(map[string]string{"data": "Hello from Mux & mongoDB"})
	}).Methods("GET") // Specifying that this route only responds to HTTP GET requests

	// Starting the HTTP server on port 6000 with the configured router
	log.Fatal(http.ListenAndServe(":6000", router))
}
