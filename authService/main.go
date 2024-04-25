package main

import (
	"auth-svc/adapters/delivery"
	"auth-svc/internal/repository"
	"auth-svc/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authRepo := repository.NewAuthRepository()
	authSvc := service.NewAuthService(authRepo)

	authHandler := delivery.NewAuthHandler(authSvc)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))

	e.POST("/login", authHandler.UserLoginHandler)
	e.GET("/revoke-token", authHandler.RevokeTokenHandler)

	// 	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}
