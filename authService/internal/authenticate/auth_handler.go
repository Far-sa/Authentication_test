package authenticate

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// User represents a user entity
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// AuthHandler handles user authentication requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get user credentials
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate authentication logic
	// In a real scenario, you would verify the user's credentials against a database or external system
	// Here, we'll just check if the username and password are "admin"
	if user.Username != "admin" || user.Password != "admin" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token for the authenticated user
	token, err := generateJWTToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Return JWT token in the response
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateJWTToken generates a JWT token for the given username
func generateJWTToken(username string) (string, error) {
	// Set token expiration time to 24 hours
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create JWT claims
	claims := JWTClaims{
		Username: username,
		Exp:      expirationTime,
	}

	// Encode claims to JSON
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// In a real scenario, you would sign the JWT token using a private key
	// For simplicity, we'll just return a base64-encoded string representation of the claims
	return string(claimsJSON), nil
}

// func main() {
// 	// HTTP route for user authentication
// 	http.HandleFunc("/authenticate", AuthHandler)

// 	// Start HTTP server
// 	fmt.Println("Authentication service is running on :8082...")
// 	log.Fatal(http.ListenAndServe(":8082", nil))
// }

// ! handle token revocation requests
// Handler for revoking tokens
func RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse token identifier from request
	tokenID := r.URL.Query().Get("token_id")

	// Add token to the list of revoked tokens in the database
	if err := db.AddRevokedToken(tokenID); err != nil {
		// Handle error
		http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

// Middleware for token validation in Traefik
func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from request
		token := ExtractTokenFromRequest(r)

		// Validate token
		if isValid := ValidateToken(token); !isValid {
			// Token is invalid or revoked
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to next handler
		next.ServeHTTP(w, r)
	})
}

// ExtractTokenFromRequest extracts JWT token from request
func ExtractTokenFromRequest(r *http.Request) string {
	// Extract token from request headers, cookies, or query parameters
	// Example: Authorization: Bearer <token>
	token := r.Header.Get("Authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}

	// Extract token from cookies or query parameters if needed

	return ""
}

// ValidateToken validates JWT token against the list of revoked tokens
func ValidateToken(token string) bool {
	// Check if token is revoked
	if db.IsTokenRevoked(token) {
		return false
	}

	// Validate token signature, expiration, etc.
	// Example: Use JWT library to parse and validate token

	return true
}
