package scheduler

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/johnroshan2255/core-service/internal/config"
	"github.com/johnroshan2255/core-service/internal/document/models"
	"github.com/johnroshan2255/core-service/internal/document/service"
	notificationv1 "github.com/johnroshan2255/core-service/proto/notification/v1"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ExpiryScheduler struct {
	documentService *service.Service
	conn            *grpc.ClientConn
	client          notificationv1.NotificationServiceClient
	serviceKey      string
	daysBeforeExpiry int
	cronScheduler   *cron.Cron
}

func NewExpiryScheduler(docService *service.Service, cfg *config.Config, daysBeforeExpiry int) *ExpiryScheduler {
	var conn *grpc.ClientConn
	var client notificationv1.NotificationServiceClient

	if cfg.CoreNotificationServiceAddr != "" {
		var creds credentials.TransportCredentials
		if cfg.TLSEnabled {
			config := &tls.Config{
				InsecureSkipVerify: false,
			}
			creds = credentials.NewTLS(config)
			log.Printf("ExpiryScheduler: Connecting to notification service with TLS at %s", cfg.CoreNotificationServiceAddr)
		} else {
			creds = insecure.NewCredentials()
			log.Printf("ExpiryScheduler: WARNING: Connecting to notification service without TLS at %s", cfg.CoreNotificationServiceAddr)
		}

		var err error
		conn, err = grpc.NewClient(cfg.CoreNotificationServiceAddr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Printf("ExpiryScheduler: Failed to connect to notification service: %v", err)
		} else {
			client = notificationv1.NewNotificationServiceClient(conn)
			log.Printf("ExpiryScheduler: Connected to notification service at %s", cfg.CoreNotificationServiceAddr)
		}
	} else {
		log.Printf("ExpiryScheduler: CoreNotificationServiceAddr not set. Notifications will be disabled.")
	}

	return &ExpiryScheduler{
		documentService:    docService,
		conn:                conn,
		client:              client,
		serviceKey:          cfg.ServiceKey,
		daysBeforeExpiry:   daysBeforeExpiry,
		cronScheduler:      cron.New(cron.WithSeconds()),
	}
}

func (s *ExpiryScheduler) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *ExpiryScheduler) Start(ctx context.Context) {
	schedule := "0 0 9 * * *"
	
	_, err := s.cronScheduler.AddFunc(schedule, func() {
		s.checkExpiringDocuments(ctx)
	})
	
	if err != nil {
		log.Printf("ExpiryScheduler: Failed to schedule job: %v", err)
		return
	}
	
	s.cronScheduler.Start()
	log.Printf("ExpiryScheduler: Started checking for expiring documents daily at 9:00 AM")
}

func (s *ExpiryScheduler) Stop() {
	s.cronScheduler.Stop()
	log.Printf("ExpiryScheduler: Stopped")
}

func (s *ExpiryScheduler) checkExpiringDocuments(ctx context.Context) {
	log.Printf("ExpiryScheduler: Checking for expiring documents...")
	
	docs, err := s.documentService.GetExpiringDocuments(ctx, s.daysBeforeExpiry)
	if err != nil {
		log.Printf("ExpiryScheduler: Failed to get expiring documents: %v", err)
		return
	}
	
	log.Printf("ExpiryScheduler: Found %d expiring documents", len(docs))
	
	for _, doc := range docs {
		if err := s.sendExpiryNotification(ctx, doc); err != nil {
			log.Printf("ExpiryScheduler: Failed to send notification for document %d: %v", doc.ID, err)
			continue
		}
		
		if err := s.documentService.MarkNotificationSent(ctx, doc.ID); err != nil {
			log.Printf("ExpiryScheduler: Failed to mark notification sent for document %d: %v", doc.ID, err)
		}
	}
	
	expiredDocs, err := s.documentService.GetExpiredDocuments(ctx)
	if err != nil {
		log.Printf("ExpiryScheduler: Failed to get expired documents: %v", err)
		return
	}
	
	log.Printf("ExpiryScheduler: Found %d expired documents", len(expiredDocs))
	
	for _, doc := range expiredDocs {
		if err := s.sendExpiredNotification(ctx, doc); err != nil {
			log.Printf("ExpiryScheduler: Failed to send expired notification for document %d: %v", doc.ID, err)
		}
	}
}

func (s *ExpiryScheduler) createContextWithAuth(ctx context.Context) context.Context {
	if s.serviceKey == "" {
		return ctx
	}
	md := metadata.New(map[string]string{
		"service-key": s.serviceKey,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func (s *ExpiryScheduler) sendExpiryNotification(ctx context.Context, doc *models.Document) error {
	if s.client == nil {
		log.Printf("ExpiryScheduler: Notification client not available, skipping notification for document %d", doc.ID)
		return nil
	}
	
	daysUntilExpiry := 0
	if doc.ExpiryDate != nil {
		daysUntilExpiry = int(time.Until(*doc.ExpiryDate).Hours() / 24)
	}
	
	userEmail := s.getUserEmailFromDocument(doc)
	if userEmail == "" {
		return fmt.Errorf("user email not found for document")
	}
	
	expiryDateStr := ""
	if doc.ExpiryDate != nil {
		expiryDateStr = doc.ExpiryDate.Format(time.RFC3339)
	}
	
	message := fmt.Sprintf(
		"Your document '%s' (Category: %s) will expire in %d days on %s. Please renew it soon.",
		doc.Name,
		doc.Category,
		daysUntilExpiry,
		doc.ExpiryDate.Format("2006-01-02"),
	)
	
	ctx = s.createContextWithAuth(ctx)
	req := &notificationv1.DocumentExpiryRequest{
		UserUuid:        doc.UserUUID,
		Email:           userEmail,
		DocumentName:   doc.Name,
		DocumentCategory: string(doc.Category),
		ExpiryDate:     expiryDateStr,
		DaysUntilExpiry: int32(daysUntilExpiry),
		IsExpired:      false,
		Message:        message,
	}
	
	_, err := s.client.NotifyDocumentExpiry(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	
	log.Printf("ExpiryScheduler: Sent expiry notification for document %d to user %s", doc.ID, doc.UserUUID)
	return nil
}

func (s *ExpiryScheduler) sendExpiredNotification(ctx context.Context, doc *models.Document) error {
	if s.client == nil {
		log.Printf("ExpiryScheduler: Notification client not available, skipping notification for document %d", doc.ID)
		return nil
	}
	
	userEmail := s.getUserEmailFromDocument(doc)
	if userEmail == "" {
		return fmt.Errorf("user email not found for document")
	}
	
	expiryDateStr := ""
	daysUntilExpiry := 0
	if doc.ExpiryDate != nil {
		expiryDateStr = doc.ExpiryDate.Format(time.RFC3339)
		daysUntilExpiry = int(time.Until(*doc.ExpiryDate).Hours() / 24)
	}
	
	message := fmt.Sprintf(
		"Your document '%s' (Category: %s) has expired on %s. Please renew it immediately.",
		doc.Name,
		doc.Category,
		doc.ExpiryDate.Format("2006-01-02"),
	)
	
	ctx = s.createContextWithAuth(ctx)
	req := &notificationv1.DocumentExpiryRequest{
		UserUuid:        doc.UserUUID,
		Email:           userEmail,
		DocumentName:   doc.Name,
		DocumentCategory: string(doc.Category),
		ExpiryDate:     expiryDateStr,
		DaysUntilExpiry: int32(daysUntilExpiry),
		IsExpired:      true,
		Message:        message,
	}
	
	_, err := s.client.NotifyDocumentExpiry(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	
	log.Printf("ExpiryScheduler: Sent expired notification for document %d to user %s", doc.ID, doc.UserUUID)
	return nil
}

func (s *ExpiryScheduler) getUserEmailFromDocument(doc *models.Document) string {
	return fmt.Sprintf("user-%s@example.com", doc.UserUUID)
}
