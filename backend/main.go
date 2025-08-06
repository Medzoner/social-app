package main

import (
	"context"
	"log"

	"social-app/internal/wire"
)

func main() {
	srv, err := wire.InitServer(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := srv.Serve(); err != nil {
		log.Fatalf("Error server: %v", err)
	}
}
