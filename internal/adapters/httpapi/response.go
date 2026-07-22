package httpapi

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a small helper to avoid repeating header/status/encode
// boilerplate in every handler.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
