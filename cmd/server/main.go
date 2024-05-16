package main

import (
	"log"
	"net/http"
	"zpe-cloud-user-management-service/config"
	"zpe-cloud-user-management-service/internal/user"
)

func main() {
	cfg := config.LoadConfig()

	user.InitializeStorage()

	mux := http.NewServeMux()

	setupRoutes(mux)

	log.Printf("Server running on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}

func setupRoutes(mux *http.ServeMux) {
	mux.Handle("/users", http.HandlerFunc(user.HandleUsers))
	mux.Handle("/users/", http.HandlerFunc(user.HandleUser))
	mux.Handle("/users/roles/", http.HandlerFunc(user.HandleUserRoles))
}
