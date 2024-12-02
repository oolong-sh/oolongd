package api

import (
	"fmt"
	"net/http"
	"os"
	"slices"
)

// TODO: whitelist native app and web urls
var allowedOrigins = []string{
	"http://localhost:8000",
	"http://localhost:5173",
}

func checkOrigin(w http.ResponseWriter, r *http.Request) error {
	origin := r.Header.Get("Origin")

	// check if the origin is whitelisted
	if origin != "" {
		if !slices.Contains(allowedOrigins, origin) {
			http.Error(w, "Request origin not in allow list", http.StatusForbidden)
			return fmt.Errorf("Requesting client not in allow list. Origin: %s\n", origin)
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	return nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
