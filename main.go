package main

import (
	"fmt"
	"gokatan/roots/configuration"
	"log"
	"net/http"
)

func main() {

    configuration.NewApplicationBuilder().WithKernel().Boot();

    // Define a basic handler function
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Go-Rose! You requested: %s", r.URL.Path)
    })

    // Start the server
    port := ":8000"
    fmt.Printf("Starting server on port %s\n", port)
    err := http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}

