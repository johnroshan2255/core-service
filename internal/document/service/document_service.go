package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/johnroshan2255/core-service/internal/document/models"
	"github.com/johnroshan2255/core-service/internal/document/repos"
	"gorm.io/gorm"
)

type Service struct {
	repo repos.Repository
}

func NewService(repo repos.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ValidateDocument(fileName, mimeType string, fileSize int64) error {
	ext := strings.ToLower(filepath.Ext(fileName))
	
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"}
	allowedMimes := []string{"image/jpeg", "image/png", "image/gif", "application/pdf"}
	
	extValid := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			extValid = true
			break
		}
	}
	
	if !extValid {
		return fmt.Errorf("invalid file type. Allowed types: jpg, jpeg, png, gif, pdf")
	}
	
	mimeValid := false
	for _, allowedMime := range allowedMimes {
		if mimeType == allowedMime {
			mimeValid = true
			break
		}
	}
	
	if !mimeValid {
		return fmt.Errorf("invalid MIME type: %s", mimeType)
	}
	
	maxSize := int64(10 * 1024 * 1024)
	if fileSize > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 10MB")
	}
	
	return nil
}

func (s *Service) DetermineDocumentType(mimeType string) models.DocumentType {
	if strings.HasPrefix(mimeType, "image/") {
		return models.DocumentTypeImage
	}
	return models.DocumentTypePDF
}

func (s *Service) ExtractDataFromDocument(ctx context.Context, doc *models.Document) (map[string]interface{}, error) {
	extractedData := make(map[string]interface{})
	
	extractedData["file_name"] = doc.FileName
	extractedData["file_type"] = string(doc.Type)
	extractedData["file_size"] = doc.FileSize
	extractedData["mime_type"] = doc.MimeType
	
	if doc.IssueDate != nil {
		extractedData["issue_date"] = doc.IssueDate.Format("2006-01-02")
	}
	
	if doc.ExpiryDate != nil {
		extractedData["expiry_date"] = doc.ExpiryDate.Format("2006-01-02")
		daysUntilExpiry := int(time.Until(*doc.ExpiryDate).Hours() / 24)
		extractedData["days_until_expiry"] = daysUntilExpiry
		
		if daysUntilExpiry < 0 {
			extractedData["status"] = "expired"
		} else if daysUntilExpiry <= 30 {
			extractedData["status"] = "expiring_soon"
		} else {
			extractedData["status"] = "active"
		}
	}
	
	extractedData["category"] = string(doc.Category)
	extractedData["extracted_at"] = time.Now().Format(time.RFC3339)
	
	return extractedData, nil
}

func (s *Service) CreateDocument(ctx context.Context, userUUID string, doc *models.Document) (*models.Document, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}
	
	if doc.Name == "" {
		return nil, fmt.Errorf("document name is required")
	}
	
	if doc.FilePath == "" {
		return nil, fmt.Errorf("file path is required")
	}
	
	doc.UserUUID = userUUID
	doc.Status = models.DocumentStatusActive
	
	if doc.ExpiryDate != nil {
		if time.Now().After(*doc.ExpiryDate) {
			doc.Status = models.DocumentStatusExpired
		} else {
			daysUntilExpiry := int(time.Until(*doc.ExpiryDate).Hours() / 24)
			if daysUntilExpiry <= 30 {
				doc.Status = models.DocumentStatusExpiring
			}
		}
	}
	
	extractedData, err := s.ExtractDataFromDocument(ctx, doc)
	if err != nil {
		log.Printf("DocumentService: Failed to extract data: %v", err)
	} else {
		extractedJSON, _ := json.Marshal(extractedData)
		doc.ExtractedData = string(extractedJSON)
	}
	
	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	
	log.Printf("DocumentService: Created document %d for user %s", doc.ID, userUUID)
	return doc, nil
}

func (s *Service) GetDocument(ctx context.Context, userUUID string, id uint) (*models.Document, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}
	
	doc, err := s.repo.GetByUUID(ctx, userUUID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	
	return doc, nil
}

func (s *Service) GetUserDocuments(ctx context.Context, userUUID string, limit, offset int) ([]models.Document, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}
	
	docs, err := s.repo.GetByUserUUID(ctx, userUUID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	
	return docs, nil
}

func (s *Service) UpdateDocument(ctx context.Context, userUUID string, id uint, updates map[string]interface{}) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}
	
	doc, err := s.repo.GetByUUID(ctx, userUUID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to get document: %w", err)
	}
	
	if name, ok := updates["name"].(string); ok && name != "" {
		doc.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		doc.Description = description
	}
	if category, ok := updates["category"].(string); ok {
		doc.Category = models.DocumentCategory(category)
	}
	if expiryDate, ok := updates["expiry_date"].(*time.Time); ok && expiryDate != nil {
		doc.ExpiryDate = expiryDate
		if time.Now().After(*expiryDate) {
			doc.Status = models.DocumentStatusExpired
		} else {
			daysUntilExpiry := int(time.Until(*expiryDate).Hours() / 24)
			if daysUntilExpiry <= 30 {
				doc.Status = models.DocumentStatusExpiring
			} else {
				doc.Status = models.DocumentStatusActive
			}
		}
	}
	
	if err := s.repo.Update(ctx, doc); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	
	log.Printf("DocumentService: Updated document %d for user %s", id, userUUID)
	return nil
}

func (s *Service) DeleteDocument(ctx context.Context, userUUID string, id uint) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}
	
	_, err := s.repo.GetByUUID(ctx, userUUID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to get document: %w", err)
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	log.Printf("DocumentService: Deleted document %d for user %s", id, userUUID)
	return nil
}

func (s *Service) GetExpiringDocuments(ctx context.Context, daysBeforeExpiry int) ([]models.Document, error) {
	docs, err := s.repo.GetExpiringDocuments(ctx, daysBeforeExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring documents: %w", err)
	}
	return docs, nil
}

func (s *Service) GetExpiredDocuments(ctx context.Context) ([]models.Document, error) {
	docs, err := s.repo.GetExpiredDocuments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired documents: %w", err)
	}
	return docs, nil
}

func (s *Service) MarkNotificationSent(ctx context.Context, id uint) error {
	return s.repo.UpdateNotificationSent(ctx, id, true)
}

