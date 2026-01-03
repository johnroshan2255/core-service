package notification

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/johnroshan2255/core-service/internal/notification"
	notificationv1 "github.com/johnroshan2255/core-service/proto/notification/v1"
)

// Handler implements the gRPC notification service
type Handler struct {
	notificationv1.UnimplementedNotificationServiceServer
	service *notification.NotificationService
}

// NewHandler creates a new notification gRPC handler
func NewHandler(notificationService *notification.NotificationService) *Handler {
	return &Handler{
		service: notificationService,
	}
}

// NotifyUserCreated handles the gRPC call for user creation notifications
func (h *Handler) NotifyUserCreated(ctx context.Context, req *notificationv1.UserCreatedRequest) (*notificationv1.UserCreatedResponse, error) {
	if req.UserUuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_uuid is required")
	}
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	log.Printf("NotificationHandler: Received NotifyUserCreated request - UUID: %s, Email: %s, Username: %s",
		req.UserUuid, req.Email, req.Username)

	err := h.service.NotifyUserCreated(ctx, req.UserUuid, req.Email, req.Username)
	if err != nil {
		log.Printf("NotificationHandler: Error processing user created notification: %v", err)
		return &notificationv1.UserCreatedResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send notification: %v", err),
		}, status.Errorf(codes.Internal, "failed to process notification: %v", err)
	}

	log.Printf("NotificationHandler: Successfully processed user created notification for %s", req.Email)
	return &notificationv1.UserCreatedResponse{
		Success: true,
		Message: "Notification sent successfully",
	}, nil
}

