package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"chirpy/internal/database"
	"chirpy/internal/handlers"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Load env file from current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error laoding env: %v", err)
	}

	// Read the db variable
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}

	// Get JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	log.Printf("Connecting to DB at: %s", dbURL)

	// Connect to db
	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	defer dbConnection.Close()

	// Initiate queries
	dbQueries := database.New(dbConnection)

	// Create config struct
	apiCfg := handlers.NewAPIConfig(dbQueries, platform, jwtSecret, polkaKey)
	mux := http.NewServeMux() // routes requests

	// Raw file server handler
	fileServer := http.FileServer(http.Dir("."))
	fileServerHandler := http.StripPrefix("/app/", fileServer)

	wrappedFileServer := apiCfg.MiddlewareHitCounter(fileServerHandler)

	// Register wrapped handlers to app
	mux.Handle("/app/", wrappedFileServer)
	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerGetChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlerPolkaWebhook)

	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUpdateUser)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandlerDeleteChirp)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server running on http://localhost:8080...")
	log.Fatal(server.ListenAndServe())
}
