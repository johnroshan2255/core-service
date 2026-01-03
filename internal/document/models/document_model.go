package models

import (
	"time"

	"gorm.io/gorm"
)

type DocumentType string

const (
	DocumentTypeImage DocumentType = "image"
	DocumentTypePDF   DocumentType = "pdf"
)

type DocumentStatus string

const (
	DocumentStatusActive   DocumentStatus = "active"
	DocumentStatusExpired  DocumentStatus = "expired"
	DocumentStatusExpiring DocumentStatus = "expiring"
)

type DocumentCategory string

const (
	DocumentCategoryWarranty      DocumentCategory = "warranty"
	DocumentCategoryPollutionCert DocumentCategory = "pollution_certificate"
	DocumentCategoryInsurance     DocumentCategory = "insurance"
	DocumentCategoryLicense       DocumentCategory = "license"
	DocumentCategoryOther         DocumentCategory = "other"
)

type Document struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserUUID    string         `gorm:"type:uuid;index;not null" json:"user_uuid"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    DocumentCategory `gorm:"type:varchar(50);not null" json:"category"`
	Type        DocumentType   `gorm:"type:varchar(20);not null" json:"type"`
	FileName    string         `gorm:"type:varchar(500);not null" json:"file_name"`
	FilePath    string         `gorm:"type:varchar(1000);not null" json:"file_path"`
	FileSize    int64          `gorm:"type:bigint" json:"file_size"`
	MimeType    string         `gorm:"type:varchar(100)" json:"mime_type"`
	IssueDate   *time.Time     `json:"issue_date,omitempty"`
	ExpiryDate  *time.Time     `gorm:"index" json:"expiry_date,omitempty"`
	Status      DocumentStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	Metadata    string         `gorm:"type:jsonb" json:"metadata,omitempty"`
	ExtractedData string        `gorm:"type:jsonb" json:"extracted_data,omitempty"`
	NotificationSent bool       `gorm:"default:false" json:"notification_sent"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

