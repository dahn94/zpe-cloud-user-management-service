package main

import (
	"log"
	"net/http"
	"zpe-cloud-user-management-service/config"
	"zpe-cloud-user-management-service/internal/app"
)

func main() {
	cfg := config.LoadConfig()

	app.InitializeStorage()

	mux := http.NewServeMux()

	setupRoutes(mux)

	log.Printf("Server running on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}

func setupRoutes(mux *http.ServeMux) {
	mux.Handle("/users", http.HandlerFunc(app.HandleUsers))
	mux.Handle("/users/", http.HandlerFunc(app.HandleUser))
	mux.Handle("/users/roles/", http.HandlerFunc(app.HandleUserRoles))
}
