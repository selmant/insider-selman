package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type WebhookPayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func basicAuthMiddleware(username, password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u, p, ok := c.Request().BasicAuth()
			if !ok || u != username || p != password {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	}
}

func handleWebhook(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return echo.NewHTTPError(http.StatusMethodNotAllowed, "Only POST method is allowed")
	}

	_, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var payload WebhookPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	log.Printf("Received webhook: %+v\n", payload)

	messageID := uuid.New().String()

	response := WebhookResponse{
		Message:   "Accepted",
		MessageID: messageID,
	}
	log.Printf("Responding with: %+v\n", response)

	return c.JSON(http.StatusAccepted, response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Use(basicAuthMiddleware(os.Getenv("WEBHOOK_USERNAME"), os.Getenv("WEBHOOK_PASSWORD")))

	e.POST("/message", handleWebhook)

	fmt.Println("Starting webhook simulator on :8080")

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
