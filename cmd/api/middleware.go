package main

import (
	"GoSocial/internal/auth"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// contextKey is a custom type to avoid context key collisions
type contextKey string

const userContextKey contextKey = "user"

func getAuthorizationToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}
	headerParts := strings.SplitN(authHeader, " ", 2)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}
	if headerParts[1] == "" {
		return "", fmt.Errorf("token missing from authorization header")
	}
	return headerParts[1], nil
}

// authenticate middleware extracts and validates JWT token from Authorization header
func (app *application) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token using helper
		token, err := getAuthorizationToken(r)
		if err != nil {
			app.unauthorizedError(w, err.Error())
			return
		}
		// Validate the token using our JWT package
		claims, err := auth.ValidateToken(token, app.config.jwt.secret)
		if err != nil {
			app.unauthorizedError(w, fmt.Sprintf("invalid token: %v", err))
			return
		}

		// Add user info to context so handlers can access it
		ctx := context.WithValue(r.Context(), userContextKey, claims)

		// Continue to the next handler with the user in context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// requireAuth wraps a handler to require authentication
func (app *application) requireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return app.authenticate(handler)
}

// requireAdmin wraps a handler to require admin privileges (placeholder for future)
func (app *application) requireAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return app.authenticate(func(w http.ResponseWriter, r *http.Request) {
		// For now, just check if user is authenticated
		// In the future, you could check user roles/permissions here
		claims := app.getUserFromContext(r)
		if claims == nil {
			app.unauthorizedError(w, "authentication required")
			return
		}

		// TODO: Add admin role checking logic here when you have user roles
		// For now, any authenticated user can access admin endpoints

		handler.ServeHTTP(w, r)
	})
}

// Helper function to get user claims from request context
func (app *application) getUserFromContext(r *http.Request) *auth.Claims {
	user, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return user
}
