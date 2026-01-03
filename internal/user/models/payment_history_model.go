package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentHistory struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserUUID      string         `gorm:"type:uuid;index;not null" json:"user_uuid"`
	TransactionID string         `gorm:"type:varchar(255);uniqueIndex" json:"transaction_id"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency      string         `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Status        string         `gorm:"type:varchar(50);not null" json:"status"`
	PaymentMethod string         `gorm:"type:varchar(50)" json:"payment_method"`
	Description   string         `gorm:"type:text" json:"description"`
	InvoiceID     string         `gorm:"type:uuid" json:"invoice_id,omitempty"`
	PaidAt        *time.Time     `json:"paid_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

