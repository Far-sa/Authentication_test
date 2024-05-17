package main

import (
	"auth-svc/adapters/config"
	"auth-svc/adapters/delivery"
	"auth-svc/adapters/messaging"
	"auth-svc/adapters/repository/db"
	"auth-svc/adapters/repository/migrator"
	"auth-svc/adapters/repository/postgres"
	"auth-svc/internal/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.ReadInConfig()
}

func main() {

	configAdapter, err := config.NewViperAdapter()
	if err != nil {
		fmt.Println("failed to load configuration", err)
	}

	dbPool, err := db.GetConnectionPool(configAdapter) // Use dedicated function (if using db package)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close() // Close the pool when done (consider connection pool management)

	authRepo := postgres.NewAuthRepository(dbPool)

	mgr, _ := migrator.NewMigrator(dbPool, "database/migrations")
	mgr.MigrateUp()

	log.Println("Migrations completed successfully!")

	//connectionString := "amqp://guest:guest@localhost:5672/"
	eventPublisher, err := messaging.NewRabbitMQClient(configAdapter)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	authSvc := service.NewAuthService(configAdapter, authRepo, eventPublisher)

	authHandler := delivery.NewAuthHandler(authSvc)

	// Initialize Echo
	e := echo.New()
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.POST("/login", authHandler.Login)
	// e.GET("/revoke-token", authHandler.RevokeToken)

	// 	// Start server
	// e.Logger.Fatal(e.Start(":5001"))

	//! Shutdown Gracefully
	// Start server in a separate goroutine
	serverStopped := make(chan error, 1)
	go func() {
		err := e.Start(":5001")
		if err != nil {
			e.Logger.Fatal(err) // Handle server startup error critically
		}
		serverStopped <- err // Signal server shutdown (if any error)
	}()

	// Wait for interrupt signal or server error to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("Received interrupt signal, initiating graceful shutdown...")
	case err := <-serverStopped:
		if err != nil {
			log.Printf("Server exited with error: %v\n", err)
		}
		log.Println("Server stopped, initiating graceful shutdown...")
	}

	// Perform cleanup tasks here before exiting
	log.Println("Shutting down server...")

	// Close resources with error handling
	if err := dbPool.Close(); err != nil {
		log.Println("Error closing database pool:", err)
	}
	if err := eventPublisher.Close(); err != nil {
		log.Println("Error closing event publisher:", err)
	}

	log.Println("Server gracefully stopped")

}
