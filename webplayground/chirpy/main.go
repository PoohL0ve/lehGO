package main

import (
	"log"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func main() {
	mux := http.NewServeMux() // routes requests
	// Convert current (.) directory to http.Dir
	mux.HandleFunc("/healthz", healthHandler) // register readiness endpoint
	fileDir := http.Dir(".")
	// Handler pointing to directory
	fileServer := http.FileServer(fileDir)

	// register server to the root path
	mux.Handle("/app/", http.StripPrefix("/app/", fileServer))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server running on http://localhost:8080...")
	log.Fatal(server.ListenAndServe())
}
