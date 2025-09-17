package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON sends a JSON response
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// writeError sends a JSON error response
func (app *application) writeError(w http.ResponseWriter, status int, message string) {
	errorResponse := map[string]string{"error": message}
	// If writeJSON fails, fall back to plain text
	if err := app.writeJSON(w, status, errorResponse); err != nil {
		http.Error(w, message, status)
	}
}

// readJSON reads JSON from request body into destination
func (app *application) readJSON(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

// Standard error responses for common cases
func (app *application) badRequestError(w http.ResponseWriter, message string) {
	app.writeError(w, http.StatusBadRequest, message)
}

func (app *application) unauthorizedError(w http.ResponseWriter, message string) {
	app.writeError(w, http.StatusUnauthorized, message)
}

func (app *application) internalServerError(w http.ResponseWriter, message string) {
	app.writeError(w, http.StatusInternalServerError, message)
}

func (app *application) conflictError(w http.ResponseWriter, message string) {
	app.writeError(w, http.StatusConflict, message)
}
