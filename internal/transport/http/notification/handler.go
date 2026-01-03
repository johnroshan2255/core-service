package notification

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/notification"
)

// Handler handles HTTP requests for notifications
type Handler struct {
	service *notification.NotificationService
}

// NewHandler creates a new HTTP notification handler
func NewHandler(service *notification.NotificationService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleUserCreated handles HTTP POST request for user creation notification
func (h *Handler) HandleUserCreated(c *gin.Context) {
	var req struct {
		UserUUID string `json:"user_uuid" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.service.NotifyUserCreated(ctx, req.UserUUID, req.Email, req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send notification",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification sent successfully",
	})
}

