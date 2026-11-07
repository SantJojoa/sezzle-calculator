// Command server runs the calculator HTTP API.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/santiagojojoa/sizzle-calculator/backend/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origin := os.Getenv("ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}

	handler := api.Logger(api.CORS(origin, api.NewRouter()))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server listening on %s (allowed origin: %s)", srv.Addr, origin)
	log.Fatal(srv.ListenAndServe())
}
