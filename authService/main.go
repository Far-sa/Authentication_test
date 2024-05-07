package main

import (
	"auth-svc/adapters/delivery"
	"auth-svc/adapters/messaging"
	"auth-svc/adapters/repository/db"
	"auth-svc/adapters/repository/migrator"
	"auth-svc/adapters/repository/postgres"
	"auth-svc/internal/service"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	//* load config
	dbConfig := db.Config{
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Port:     viper.GetInt("db.port"),
		Host:     viper.GetString("db.host"),
		DbName:   viper.GetString("db.dbName"),
	}

	rabbitConf := messaging.RabbitMQConfig{
		Host:     viper.GetString("rabbitmq.host"),
		User:     viper.GetString("rabbitmq.user"),
		Password: viper.GetString("rabbitmq.password"),
		Port:     viper.GetString("rabbitmq.port"),
	}

	dbPool, err := db.GetConnectionPool(dbConfig) // Use dedicated function (if using db package)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close() // Close the pool when done (consider connection pool management)

	authRepo := postgres.NewAuthRepository(dbPool)

	mgr := migrator.NewMigrator(dbPool, "database/migrations")
	mgr.MigrateUp()

	log.Println("Migrations completed successfully!")

	//connectionString := "amqp://guest:guest@localhost:5672/"
	eventPublisher, err := messaging.NewRabbitMQClient(rabbitConf)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	authSvc := service.NewAuthService(authRepo, eventPublisher)

	authHandler := delivery.NewAuthHandler(authSvc)

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", authHandler.Login)
	// e.GET("/revoke-token", authHandler.RevokeToken)

	// 	// Start server
	e.Logger.Fatal(e.Start(":5001"))

}
