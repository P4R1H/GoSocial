package main

import (
	"GoSocial/internal/auth"
	"GoSocial/internal/store"
	"context"
	"net/http"
	"strings"
	"time"
)

// Request/Response structures
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  store.User `json:"user"`
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// registerHandler handles user registration
func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Parse JSON request
	if err := app.readJSON(r, &req); err != nil {
		app.badRequestError(w, "Invalid JSON")
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		app.badRequestError(w, "username, email, and password are required")
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		app.internalServerError(w, "Failed to process password")
		return
	}

	// Create user in database with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := app.store.Users.Create(ctx, user); err != nil {
		// Simple duplicate check - you might want to improve this
		if contains(err.Error(), "duplicate") || contains(err.Error(), "unique") {
			app.conflictError(w, "Username or email already exists")
			return
		}
		app.internalServerError(w, "Failed to create user")
		return
	}

	// Generate JWT token using config
	token, err := auth.GenerateToken(
		user.ID,
		user.Username,
		app.config.jwt.secret, // ← From your config
		app.config.jwt.issuer, // ← From your config
		app.config.jwt.expiry, // ← From your config
	)
	if err != nil {
		app.internalServerError(w, "Failed to generate token")
		return
	}

	// Return response (password is hidden by json:"-" tag)
	user.Password = ""
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	app.writeJSON(w, http.StatusCreated, response)
}

// loginHandler handles user login
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse JSON request
	if err := app.readJSON(r, &req); err != nil {
		app.badRequestError(w, "Invalid JSON")
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		app.badRequestError(w, "Email and password are required")
		return
	}

	// Get user from database by email with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := app.store.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal whether user exists or not
		app.unauthorizedError(w, "Invalid credentials")
		return
	}

	// Verify password
	if err := auth.ComparePassword(user.Password, req.Password); err != nil {
		app.unauthorizedError(w, "Invalid credentials")
		return
	}

	// Generate JWT token using config
	token, err := auth.GenerateToken(
		user.ID,
		user.Username,
		app.config.jwt.secret, // ← From your config
		app.config.jwt.issuer, // ← From your config
		app.config.jwt.expiry, // ← From your config
	)
	if err != nil {
		app.internalServerError(w, "Failed to generate token")
		return
	}

	// Return response (password is hidden by json:"-" tag)
	user.Password = ""
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	app.writeJSON(w, http.StatusOK, response)
}
