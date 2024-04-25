package authenticate

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// TokenHandler handles token-related HTTP requests
type authHandler struct {
	authService authService
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthHandler(authService authService) *authHandler {
	return &authHandler{authService}
}

func (h *authHandler) GenerateTokenHandler(c echo.Context) error {
	// Extract user ID and roles from request
	userID := c.FormValue("userID")
	roles := c.FormValue("roles") // Assuming roles are provided as a comma-separated string

	// Convert roles string to slice
	roleSlice := strings.Split(roles, ",")

	// Generate token
	token, err := h.authService.GenerateToken(userID, roleSlice)
	if err != nil {
		// Handle error
		return c.String(http.StatusInternalServerError, "Failed to generate token")
	}

	// Respond with generated token
	return c.String(http.StatusOK, token)
}

// RevokeTokenHandler handles requests to revoke a token
func (h *authHandler) RevokeTokenHandler(c echo.Context) error {
	// Parse token identifier from request
	tokenID := c.QueryParam("token_id")

	// Call token service to revoke token
	if err := h.authService.AddRevokedToken(tokenID); err != nil {
		// Handle error
		return c.String(http.StatusInternalServerError, "Failed to revoke token")
	}

	// Respond with success
	return c.String(http.StatusOK, "Token revoked successfully")
}

//!--- main
// func main() {
// 	// Initialize Echo instance
// 	e := echo.New()

// 	// Initialize token repository
// 	tokenRepo := auth.NewTokenRepository()

// 	// Initialize token service with repository
// 	authService := auth.NewTokenService(tokenRepo)

// 	// Initialize token handler with service
// 	tokenHandler := auth.NewTokenHandler(authService)

// 	// Middleware
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())

// 	// Routes
// e.POST("/generate-token", authHandler.GenerateTokenHandler)
// 	e.GET("/revoke-token", tokenHandler.RevokeTokenHandler)

// 	// Start server
// 	e.Logger.Fatal(e.Start(":8080"))
// }
