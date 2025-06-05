package main

import (
	"log"
	"net/http"
	"sample-go-weather-hour/internal/handlers"
)

func main() {
	h := handlers.NewHandler()

	http.HandleFunc("/api/v1/forecast", h.GetForecastHandler)
	http.HandleFunc("/status", h.GetStatusHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
