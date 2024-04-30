package main

import (
	"log"
)

func main() {
	server := NewAPIServer(":8081")
	err := server.Run()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
