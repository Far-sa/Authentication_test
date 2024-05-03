package main

import (
	"fmt"
	"log"
	"user-svc/adapters/config"
	"user-svc/adapters/delivery/httpServer"
	"user-svc/adapters/logger"
	"user-svc/adapters/messaging"
	"user-svc/adapters/repository/mysql"
	userService "user-svc/internal/service"
)

func main() {

	//? Initialize configuration adapter
	configAdapter, err := config.NewViperAdapter(".")
	if err != nil {
		fmt.Println("failed to load configuration", err)
	}

	//* Initialize Prometheus metrics adapter
	//prometheusAdapter := metrics.NewPrometheus()
	zapLogger, _ := logger.NewZapLogger(configAdapter)

	//* Initialize repositories and services
	userRepository := mysql.New(configAdapter, zapLogger)

	fmt.Println("host rabbit:", configAdapter.GetBrokerConfig().Host)
	//connectionString := "amqp://guest:guest@localhost:5672/"
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
