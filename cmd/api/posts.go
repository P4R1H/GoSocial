package main

import (
	"net/http"
)

// getPostsHandler returns posts for authenticated users
func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user claims from the context (added by our authenticate middleware)
	claims := app.getUserFromContext(r)
	if claims == nil {
		app.unauthorizedError(w, "authentication required")
		return
	}

	// Example response showing posts (in a real app, you'd query the database)
	postsResponse := map[string]interface{}{
		"message":  "Here are your posts!",
		"user_id":  claims.UserID,
		"username": claims.Username,
		"posts": []map[string]interface{}{
			{
				"id":      1,
				"title":   "My first post",
				"content": "Hello, world!",
				"author":  claims.Username,
			},
			{
				"id":      2,
				"title":   "Another post",
				"content": "This is a protected endpoint!",
				"author":  claims.Username,
			},
		},
	}

	if err := app.writeJSON(w, http.StatusOK, postsResponse); err != nil {
		app.internalServerError(w, "failed to write response")
	}
}
