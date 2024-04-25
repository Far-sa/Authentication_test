package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	HandleAccessToken(c echo.Context) error
	HandleRefreshToken(c echo.Context) error
}

type authHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) AuthHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) HandleAccessToken(c echo.Context) error {
	// Parse user ID from request context or JWT token
	userID := "user123" // Dummy user ID for demonstration

	// Generate access token
	accessToken, err := h.service.CreateAccessToken(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create access token")
	}

	return c.JSON(http.StatusOK, map[string]string{"access_token": accessToken})
}

func (h *authHandler) HandleRefreshToken(c echo.Context) error {
	// Parse user ID from request context or JWT token
	userID := "user123" // Dummy user ID for demonstration

	// Generate refresh token
	refreshToken, err := h.service.CreateRefreshToken(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create refresh token")
	}

	return c.JSON(http.StatusOK, map[string]string{"refresh_token": refreshToken})
}
