package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentDetails struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserUUID      string         `gorm:"type:uuid;uniqueIndex;not null" json:"user_uuid"`
	PaymentMethod string         `gorm:"type:varchar(50)" json:"payment_method"`
	CardLast4     string         `gorm:"type:varchar(4)" json:"card_last4"`
	CardBrand     string         `gorm:"type:varchar(50)" json:"card_brand"`
	BillingEmail  string         `gorm:"type:varchar(255)" json:"billing_email"`
	BillingName   string         `gorm:"type:varchar(255)" json:"billing_name"`
	BillingAddress string        `gorm:"type:text" json:"billing_address"`
	IsDefault     bool           `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

