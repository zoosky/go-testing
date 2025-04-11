package main

import (
	"fmt"
	"log"
	"net/http"

	_ "go-testing/docs" // Import for swagger
	"go-testing/internal/api"
	"go-testing/internal/calculator"
	"go-testing/internal/database"
)

func main() {
	// Initialize database repository
	repo := database.NewUserRepository()
	
	// Initialize calculator service
	calc := calculator.NewCalculator()
	
	// Initialize API server with dependencies
	server := api.NewServer(repo, calc)
	
	// Start server
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", server.Router()))
}