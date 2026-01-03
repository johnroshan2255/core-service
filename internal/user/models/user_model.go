package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UUID         string    `gorm:"type:uuid;uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null;column:password"`
	PhoneNumber  string    `gorm:"type:varchar(20);column:phone_number"`
	FirstName    string    `gorm:"type:varchar(100);column:first_name"`
	LastName     string    `gorm:"type:varchar(100);column:last_name"`
	TenantID     string    `gorm:"type:varchar(255);column:tenant_id"`
	Role         string    `gorm:"type:varchar(50);default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

