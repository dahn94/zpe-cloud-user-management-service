package main

import (
	"log"
	"net/http"
	"zpe-cloud-user-management-service/config"
	"zpe-cloud-user-management-service/internal/user"
)

func setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", user.HandleCreateUser)
	mux.HandleFunc("GET /users", user.HandleListUsers)
	mux.HandleFunc("GET /users/{id}", user.HandleGetUser)
	mux.HandleFunc("DELETE /users/{id}", user.HandleDeleteUser)
	mux.HandleFunc("PUT /users/roles/{id}", user.HandleUpdateUserRoles)
}

func main() {
	// Load the environment configuration
	cfg := config.LoadEnvConfig()

	// Initialize the storage for users
	user.InitializeStorage()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up the routes
	setupRoutes(mux)

	// Start the server
	log.Printf("Server running on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}
