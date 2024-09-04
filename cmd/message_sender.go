package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
	_ "insider/docs"
	"insider/external_services/database"
	"insider/external_services/message_publisher"
	"insider/internal/domains/message"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	messagePublisher := message_publisher.NewMessagePublisher()
	messageRepository := message.NewRepository(db)
	messageService := message.NewService(messageRepository, messagePublisher)

	messageHandler := message.NewHandler(messageService)

	messageHandler.RegisterRoutes(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	port := os.Getenv("MESSAGE_SENDER_PORT")

	log.Info("Starting message sender service on port ", port)
	e.Logger.Fatal(e.Start(":" + port))
}
