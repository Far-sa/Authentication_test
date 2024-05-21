package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-svc/adapters/config"
	"user-svc/adapters/delivery/httpServer"
	"user-svc/adapters/logger"
	"user-svc/adapters/messaging"
	"user-svc/adapters/metrics"
	"user-svc/adapters/repository/db"
	"user-svc/adapters/repository/mysql"
	userService "user-svc/internal/service"
)

func main() {

	//? Initialize configuration adapter
	configAdapter, err := config.NewViperAdapter()
	if err != nil {
		fmt.Println("failed to load configuration", err)
	}

	dbPool, err := db.GetConnectionPool(configAdapter) // Use dedicated function (if using db package)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close() // Close the pool when done (consider connection pool management)

	//* Initialize Prometheus metrics adapter
	prometheusAdapter := metrics.NewPrometheus()
	zapLogger, err := logger.NewZapLogger(configAdapter)
	if err != nil {
		log.Fatalf("failed to init logger-config: %v", err)
	}

	//* Initialize repositories and services
	userRepository := mysql.New(dbPool, zapLogger, prometheusAdapter)

	// mgr := migrator.New(dbPool, "infrastructure/db/migrations")
	// mgr.MigrateUp()

	// log.Println("Migrations completed successfully!")

	publisher, err := messaging.NewRabbitMQClient(configAdapter)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	userService := userService.NewService(configAdapter, userRepository, publisher, zapLogger)
	mErr := userService.StartMessageListener(context.Background())
	if mErr != nil {
		log.Fatalf("failed to start message listener: %v", mErr)
	}
	// ozzoValidator := validator.New(userRepository)

	//* Initialize grpc client
	// grpcHandler := grpcserver.New(userService)
	// grpcHandler.Start()

	//* http handler
	userHandler := httpServer.New(configAdapter, userService, zapLogger, prometheusAdapter)

	// userHandler.Serve()

	//! Shutdown gracefully
	serverStopped := make(chan error, 1)
	go func() {
		err := userHandler.Serve()
		if err != nil {
			log.Println("Server exited with error:", err)
		}
		serverStopped <- err
	}()

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
	if err := publisher.Close(); err != nil {
		log.Println("Error closing event publisher:", err)
	}

	log.Println("Server gracefully stopped")
}
