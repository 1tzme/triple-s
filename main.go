package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/internal/router"
	"triple-s/internal/structure"
	v "triple-s/internal/validator"
)

func main() {
	port, dir, help := v.InitFlags()

	if help {
		v.PrintUsage()
		return
	}

	err := v.ValidateDataDirectory(dir)
	if err != nil {
		log.Fatalf("Invalid data directory: %v", err)
	}

	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		log.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	server := structure.Server{
		Dir:  dir,
		Port: port,
	}

	mux := router.Router(&server)

	fmt.Printf("Starting server on port %s, directory %s\n", port, dir)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
