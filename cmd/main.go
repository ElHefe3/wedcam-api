package main

import (
	"fmt"
	"log"
	"net/http"
	"wedcam-api/pkg/api"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Service is running")
	})
	http.HandleFunc("/upload", api.ImageUploadHandler)

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
