package models

import (
	"time"

	"gorm.io/gorm"
)

type CompanyDetails struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserUUID    string         `gorm:"type:uuid;uniqueIndex;not null" json:"user_uuid"`
	CompanyName string         `gorm:"type:varchar(255)" json:"company_name"`
	Industry    string         `gorm:"type:varchar(100)" json:"industry"`
	Website     string         `gorm:"type:varchar(255)" json:"website"`
	Address     string         `gorm:"type:text" json:"address"`
	Phone       string         `gorm:"type:varchar(50)" json:"phone"`
	TaxID       string         `gorm:"type:varchar(100)" json:"tax_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

