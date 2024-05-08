package main

// import (
// 	"context"
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// )

// type OAuthService interface {
// 	HandleGoogleLogin(c echo.Context)
// 	HandleGoogleCallback(c echo.Context)
// }

// type oauthService struct {
// 	oauthConfig *oauth2.Config
// }

// func NewOAuthService(clientID, clientSecret, redirectURL string) OAuthService {
// 	return &oauthService{
// 		oauthConfig: &oauth2.Config{
// 			ClientID:     clientID,
// 			ClientSecret: clientSecret,
// 			RedirectURL:  redirectURL,
// 			Scopes:       []string{"openid", "profile", "email"},
// 			Endpoint:     google.Endpoint,
// 		},
// 	}
// }

// func (s *oauthService) HandleGoogleLogin(c echo.Context) error {
// 	// Redirect to OAuth provider's login page
// 	url := s.oauthConfig.AuthCodeURL("state")
// 	return c.Redirect(http.StatusFound, url)
// }

// func (s *oauthService) HandleGoogleCallback(c echo.Context) error {
// 	code := c.QueryParam("code")

// 	// Exchange authorization code for access token
// 	token, err := s.oauthConfig.Exchange(context.Background(), code)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, "Failed to exchange code for token")
// 	}

// 	// You can validate and process the token here
// 	// For simplicity, we'll just return the token in the response
// 	return c.JSON(http.StatusOK, token)
// }

// 	// Initialize OAuth service with dependencies
// 	oauth := NewOAuthService(
// 		"your-client-id",
// 		"your-client-secret",
// 		"http://localhost:8000/auth/google/callback",
// 	)

// 	// OAuth routes
// 	e.GET("/auth/google", oauth.HandleGoogleLogin)
// 	e.GET("/auth/google/callback", oauth.HandleGoogleCallback)
