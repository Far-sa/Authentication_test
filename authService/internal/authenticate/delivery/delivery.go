package delivery

import (
	"auth-svc/internal/authenticate/param"
	"auth-svc/internal/authenticate/ports"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TokenHandler handles token-related HTTP requests
type authHandler struct {
	authService ports.AuthService
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthHandler(authService ports.AuthService) authHandler {
	return authHandler{authService}
}

func (h authHandler) UserLoginHandler(c echo.Context) error {

	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	//* validate Login
	// if _, err := h.userValidator.ValidateLoginRequest(req); err != nil {
	// 	msg, code := httpmsg.Error(err)
	// 	return c.JSON(code, echo.Map{
	// 		"messsage": msg,
	// 		"errors":   err,
	// 	})
	// }

	resp, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// RevokeTokenHandler handles requests to revoke a token
func (h *authHandler) RevokeTokenHandler(c echo.Context) error {
	// Parse token identifier from request
	tokenID := c.QueryParam("token_id")

	// Call token service to revoke token
	if err := h.authService.AddRevokeToken(tokenID); err != nil {
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
// e.POST("/login", authHandler.UserLoginHandler)
// 	e.GET("/revoke-token", tokenHandler.RevokeTokenHandler)

// 	// Start server
// 	e.Logger.Fatal(e.Start(":8080"))
// }
