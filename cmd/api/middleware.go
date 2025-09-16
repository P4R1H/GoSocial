package main

// JWT Middleware Flow:
//
// 1. Client sends a request to a protected endpoint, including a JWT in the Authorization header.
//
// 2. Middleware intercepts the request before it reaches the handler.
//
// 3. Middleware extracts the JWT from the Authorization header (usually in the format "Bearer <token>").
//
// 4. Middleware verifies the JWT signature using the server's secret key or public key.
//    - If the signature is invalid or the token is malformed, respond with 401 Unauthorized.
//
// 5. Middleware checks the token's claims (e.g., expiration, issuer, audience).
//    - If the token is expired or claims are invalid, respond with 401 Unauthorized.
//
// 6. If the token is valid, middleware may extract user information (e.g., user ID, roles) from the token claims.
//
// 7. Middleware attaches the user information to the request context for use in downstream handlers.
//
// 8. Middleware passes the request to the next handler in the chain.
//
// 9. Handler can access the user information from the context to authorize actions or personalize responses.
