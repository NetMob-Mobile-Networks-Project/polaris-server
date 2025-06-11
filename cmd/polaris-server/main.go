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

	// Apply CORS middleware to the main router first
	router.Use(apiHandler.CORSMiddleware)

	// Register all routes
	apiHandler.RegisterRoutes(router)

	// Add a catch-all OPTIONS handler for any unmatched routes
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
	})

	port := config.Port

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
