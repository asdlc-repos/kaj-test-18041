package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"manager-api/handlers"
	"manager-api/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	// LEAVE_API_URL is injected by OpenChoreo at runtime via workload.yaml dependency binding.
	// It is available for potential future integration with the leave API.
	leaveAPIURL := os.Getenv("LEAVE_API_URL")
	if leaveAPIURL != "" {
		log.Printf("Leave API URL configured: %s", leaveAPIURL)
	}

	store := services.NewStore()
	h := handlers.NewHandler(store)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	addr := "0.0.0.0:" + port
	log.Printf("manager-api listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
