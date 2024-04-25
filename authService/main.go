package main

import (
	"auth-svc/adapters/delivery"
	"auth-svc/adapters/messaging"
	"auth-svc/adapters/repository"
	"auth-svc/internal/service"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authRepo := repository.NewAuthRepository()

	publisher, err := messaging.NewRabbitMQClient("puppet", "password", "localhost:5672", "users")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	authSvc := service.NewAuthService(authRepo, publisher)

	authHandler := delivery.NewAuthHandler(authSvc)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))

	e.POST("/login", authHandler.Login)
	// e.GET("/revoke-token", authHandler.RevokeToken)

	// 	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}
