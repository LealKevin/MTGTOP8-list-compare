package server

import (
	"fmt"
	"log"
	"net/http"
)

func InitServer() {
	mux := Router()

	// Serve the static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/*", fs)

	fmt.Printf("Server running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Impossible to connect to the server, error: %v", err)
		return
	}
}
