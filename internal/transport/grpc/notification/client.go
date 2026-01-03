package notification

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	notificationv1 "github.com/johnroshan2255/core-service/proto/notification/v1"
)

// ClientFactory creates gRPC clients for inter-service communication
type ClientFactory struct {
	serviceKey string
	grpcAddr   string
	useTLS     bool
}

// NewClientFactory creates a new gRPC client factory
func NewClientFactory(grpcAddr, serviceKey string, useTLS bool) *ClientFactory {
	return &ClientFactory{
		serviceKey: serviceKey,
		grpcAddr:   grpcAddr,
		useTLS:     useTLS,
	}
}

// CreateClient creates a gRPC client connection with authentication
func (f *ClientFactory) CreateClient(ctx context.Context) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials

	if f.useTLS {
		config := &tls.Config{
			InsecureSkipVerify: false,
		}
		creds = credentials.NewTLS(config)
		log.Printf("Creating gRPC client with TLS to %s", f.grpcAddr)
	} else {
		creds = insecure.NewCredentials()
		log.Printf("WARNING: Creating gRPC client without TLS to %s", f.grpcAddr)
	}

	conn, err := grpc.NewClient(f.grpcAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return conn, nil
}

// CreateContextWithAuth adds service key to context metadata for authentication
func (f *ClientFactory) CreateContextWithAuth(ctx context.Context) context.Context {
	if f.serviceKey == "" {
		return ctx
	}
	md := metadata.New(map[string]string{
		"service-key": f.serviceKey,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// Client is a gRPC client for the notification service
// This is used by other services within the same project to call notification service via gRPC
type Client struct {
	conn    *grpc.ClientConn
	factory *ClientFactory
	client  notificationv1.NotificationServiceClient
}

// NewClient creates a new notification service gRPC client
// This should be used by other services (user, invoice, document) to call notification service
func NewClient(factory *ClientFactory, ctx context.Context) (*Client, error) {
	conn, err := factory.CreateClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification client: %w", err)
	}

	client := notificationv1.NewNotificationServiceClient(conn)

	return &Client{
		conn:    conn,
		factory: factory,
		client:  client,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// NotifyUserCreated calls the notification service to send a user creation notification
func (c *Client) NotifyUserCreated(ctx context.Context, userUUID, email, username string) error {
	ctx = c.factory.CreateContextWithAuth(ctx)

	req := &notificationv1.UserCreatedRequest{
		UserUuid: userUUID,
		Email:    email,
		Username: username,
	}

	resp, err := c.client.NotifyUserCreated(ctx, req)
	if err != nil {
		log.Printf("NotificationClient: Failed to notify user created: %v", err)
		return fmt.Errorf("failed to notify user created: %w", err)
	}

	if !resp.Success {
		log.Printf("NotificationClient: Notification service returned failure: %s", resp.Message)
		return fmt.Errorf("notification service error: %s", resp.Message)
	}

	log.Printf("NotificationClient: Successfully notified user creation for %s", email)
	return nil
}

