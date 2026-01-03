package notification

import (
	"fmt"
	"log"
)

// Factory creates notification service instances
type Factory struct {
	provider Provider
}

// NewFactory creates a new notification factory
func NewFactory(providerType string) (*Factory, error) {
	var provider Provider

	switch providerType {
	case "email":
		provider = NewEmailProvider()
		log.Printf("NotificationFactory: Using email provider")
	case "mock":
		provider = NewMockProvider()
		log.Printf("NotificationFactory: Using mock provider")
	default:
		return nil, fmt.Errorf("unknown notification provider: %s", providerType)
	}

	return &Factory{
		provider: provider,
	}, nil
}

// NewService creates a new notification service with the configured provider
func (f *Factory) NewService() *NotificationService {
	return NewNotificationService(f.provider)
}

