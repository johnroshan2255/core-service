package notification

import (
	"context"
	"fmt"
	"log"
)

// NotificationService handles notification business logic
type NotificationService struct {
	provider Provider
}

// NewNotificationService creates a new notification service
func NewNotificationService(provider Provider) *NotificationService {
	return &NotificationService{
		provider: provider,
	}
}

// NotifyUserCreated sends a notification when a user is created
func (s *NotificationService) NotifyUserCreated(ctx context.Context, userUUID, email, username string) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}
	if email == "" {
		return fmt.Errorf("email is required")
	}

	log.Printf("NotificationService: Processing user created notification - UUID: %s, Email: %s, Username: %s", userUUID, email, username)

	notificationData := map[string]interface{}{
		"type":      "user_created",
		"user_uuid": userUUID,
		"email":     email,
		"username":  username,
	}

	if err := s.provider.SendNotification(ctx, email, "Welcome to our platform!", notificationData); err != nil {
		log.Printf("NotificationService: Failed to send user created notification: %v", err)
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Printf("NotificationService: Successfully sent user created notification to %s", email)
	return nil
}

func (s *NotificationService) NotifyDocumentExpiry(ctx context.Context, userUUID, email, documentName, documentCategory, expiryDate string, daysUntilExpiry int32, isExpired bool, message string) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if documentName == "" {
		return fmt.Errorf("document name is required")
	}

	log.Printf("NotificationService: Processing document expiry notification - UUID: %s, Email: %s, Document: %s, Category: %s, Expired: %v, Days: %d",
		userUUID, email, documentName, documentCategory, isExpired, daysUntilExpiry)

	notificationData := map[string]interface{}{
		"type":              "document_expiry",
		"user_uuid":         userUUID,
		"email":             email,
		"document_name":     documentName,
		"document_category": documentCategory,
		"expiry_date":       expiryDate,
		"days_until_expiry": daysUntilExpiry,
		"is_expired":        isExpired,
		"message":           message,
	}

	subject := "Document Expiry Notification"
	if isExpired {
		subject = fmt.Sprintf("Document Expired: %s", documentName)
	} else {
		subject = fmt.Sprintf("Document Expiring Soon: %s", documentName)
	}

	if err := s.provider.SendNotification(ctx, email, subject, notificationData); err != nil {
		log.Printf("NotificationService: Failed to send document expiry notification: %v", err)
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Printf("NotificationService: Successfully sent document expiry notification to %s", email)
	return nil
}
