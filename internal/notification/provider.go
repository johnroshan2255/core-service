package notification

import (
	"context"
	"log"
)

// Provider defines the interface for notification providers (email, SMS, push, etc.)
type Provider interface {
	SendNotification(ctx context.Context, recipient, subject string, data map[string]interface{}) error
}

// EmailProvider is a basic email notification provider
type EmailProvider struct {
}

// NewEmailProvider creates a new email notification provider
func NewEmailProvider() *EmailProvider {
	return &EmailProvider{}
}

// SendNotification sends an email notification
func (p *EmailProvider) SendNotification(ctx context.Context, recipient, subject string, data map[string]interface{}) error {
	log.Printf("EmailProvider: Sending email to %s with subject: %s", recipient, subject)
	log.Printf("EmailProvider: Notification data: %+v", data)
	return nil
}

// MockProvider is a mock provider for testing
type MockProvider struct{}

// NewMockProvider creates a new mock notification provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// SendNotification logs the notification without actually sending it
func (p *MockProvider) SendNotification(ctx context.Context, recipient, subject string, data map[string]interface{}) error {
	log.Printf("MockProvider: Would send notification to %s: %s - %+v", recipient, subject, data)
	return nil
}
