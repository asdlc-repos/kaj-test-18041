package main

import (
	"leave-api/handlers"
	"leave-api/services"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	leaveSvc := services.NewLeaveService()
	leaveHandler := handlers.NewLeaveHandler(leaveSvc)

	mux := http.NewServeMux()
	leaveHandler.RegisterRoutes(mux)

	addr := ":" + port
	log.Printf("leave-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
