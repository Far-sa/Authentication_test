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
	configAdapter := config.NewViperAdapter()
	err := configAdapter.LoadConfig("./config.yaml")
	if err != nil {
		fmt.Println("failed to load configuration", err)
	}

	//* Initialize Prometheus metrics adapter
	//prometheusAdapter := metrics.NewPrometheus()
	zapLogger, _ := logger.NewZapLogger(configAdapter)

	//* Initialize repositories and services
	userRepository := mysql.New(configAdapter, zapLogger)
	// TODO: add to config
	publisher, err := messaging.NewRabbitMQClient("puppet", "password", "localhost:5672")
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
