package message

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"insider/internal/models"
	"insider/internal/utils"
	"net/http"
)

// Handler struct
type Handler struct {
	service Service
}

// NewHandler creates a new Handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// QueueMessage godoc
// @Summary Queue a message for sending
// @Description Queue a message for sending
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.Message true "Message to be queued"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /messages/queue [post]
func (h *Handler) QueueMessage(c echo.Context) error {
	var message models.Message
	if err := json.NewDecoder(c.Request().Body).Decode(&message); err != nil {
		return c.JSON(http.StatusBadRequest, utils.APIResponse{Message: err.Error()})
	}

	if err := h.service.QueueMessageForSending(c.Request().Context(), message.Content, message.RecipientPhone); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.APIResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, utils.APIResponse{Message: "Message queued successfully"})
}

// GetMessages godoc
// @Summary Get messages by sent status
// @Description Get messages by sent status
// @Tags messages
// @Accept json
// @Produce json
// @Param status query string false "Sent status" Enums(pending, sent, failed)
// @Success 200 {object} GetMessagesResponse
// @Failure 400 {object} GetMessagesResponse
// @Router /messages [get]
func (h *Handler) GetMessages(c echo.Context) error {
	statusString := c.QueryParam("status")

	if statusString == "" {
		// Handle case where status is not provided
		messages, err := h.service.FindAllMessages(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.APIResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, utils.APIResponseWithData[[]models.Message]{Data: messages})
	}

	status, err := models.SentStatusFromString(statusString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.APIResponse{Message: err.Error()})
	}

	messages, err := h.service.FindMessagesBySentStatus(c.Request().Context(), status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.APIResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, utils.APIResponseWithData[[]models.Message]{Data: messages})
}

// RegisterRoutes registers the HTTP routes
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/messages/queue", h.QueueMessage)
}
