package main

import (
	"log"
	"wedcam-/pkg/server"
)

func main() {
	// Start the server
	if err := server.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
