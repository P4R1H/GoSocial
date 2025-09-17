package main

import (
	"net/http"
)

// getUserProfileHandler returns the current user's profile
// This demonstrates how to use the middleware and get user info from context
func (app *application) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user claims from the context (added by our authenticate middleware)
	claims := app.getUserFromContext(r)
	if claims == nil {
		app.unauthorizedError(w, "authentication required")
		return
	}

	// Create a response with user info from the JWT token
	userProfile := map[string]interface{}{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"message":  "This is your protected profile!",
	}

	// Send the JSON response
	if err := app.writeJSON(w, http.StatusOK, userProfile); err != nil {
		app.internalServerError(w, "failed to write response")
	}
}

// getAllUsersHandler is an example admin-only endpoint
func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the user passes admin authentication
	claims := app.getUserFromContext(r)

	adminResponse := map[string]interface{}{
		"message":     "This is an admin-only endpoint",
		"accessed_by": claims.Username,
		"admin_note":  "In the future, you can add role-based access control here",
	}

	if err := app.writeJSON(w, http.StatusOK, adminResponse); err != nil {
		app.internalServerError(w, "failed to write response")
	}
}
