package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/ascii-art", asciiArtHandler)

    fmt.Println("Server running at http://localhost:8080")

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Error starting server:", err)
    }
}