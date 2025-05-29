package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/api"
	database "github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage"
)

func main() {
	config := NewConfig()

	db, err := database.NewMySQLConnection(
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	apiHandler := api.NewHandler(
		db,
		config.JWTSecret,
	)
	apiHandler.RegisterRoutes(router)

	port := config.Port

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
