package handlers

import (
	"chirpy/internal/database"
	"fmt"
	"net/http"
	"sync/atomic"
)

// Track stateful in-memory data
type ApiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
}

// NewAPIConfig initializes a new APIConfig instance
func NewAPIConfig(db *database.Queries, platform string) *ApiConfig {
	return &ApiConfig{
		DB:       db,
		Platform: platform,
	}
}

// Intercepts requests to increment hit counter
func (cfg *ApiConfig) MiddlewareHitCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the cfg
		cfg.fileserverHits.Add(1)
		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// Prints number of hits
func (cfg *ApiConfig) PrintHits() {
	currentHits := cfg.fileserverHits.Load()
	fmt.Printf("Hits: %d\n", currentHits)
}

// Reads the atomic counter and writes the value to client as HTML
func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Safely load the current value from atomic.Int32
	hits := cfg.fileserverHits.Load()

	html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Write([]byte(html))
}

// Resets the atomic counter to 0
func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		RespondWithError(w, http.StatusForbidden, "Forbidden: Reset is only allowed in dev environment")
		return
	}

	// Safely reset atomic.Int32 to 0
	cfg.fileserverHits.Store(0)

	// Delete all users
	err := cfg.DB.DeleteUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete users")
		return
	}

	// Explicitly set content type, Go automatically defaults to the same "text/plain"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
