package main

import (
	"fmt"
	"log"
	"user-svc/adapters/config"
	"user-svc/adapters/delivery/httpServer"
	"user-svc/adapters/logger"
	"user-svc/adapters/messaging"
	"user-svc/adapters/repository/db"
	"user-svc/adapters/repository/migrator"
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
	//prometheusAdapter := metrics.NewPrometheus()
	zapLogger, err := logger.NewZapLogger(configAdapter)
	if err != nil {
		log.Fatalf("failed to init logger-config: %v", err)
	}

	//* Initialize repositories and services
	userRepository := mysql.New(dbPool, zapLogger)

	mgr := migrator.New(dbPool, "infrastructure/db/migrations")
	mgr.MigrateUp()

	log.Println("Migrations completed successfully!")

	publisher, err := messaging.NewRabbitMQClient(configAdapter)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	userService := userService.NewService(userRepository, publisher, zapLogger)

	// ozzoValidator := validator.New(userRepository)

	//* Initialize grpc client
	// grpcHandler := grpcserver.New(userService)
	// grpcHandler.Start()

	//* http handler
	userHandler := httpServer.New(configAdapter, userService, zapLogger)

	userHandler.Serve()

}
