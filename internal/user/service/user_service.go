package service

import (
	"context"
	"fmt"
	"log"

	"github.com/johnroshan2255/core-service/internal/user/models"
	"github.com/johnroshan2255/core-service/internal/user/repos"
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

func (s *Service) GetProfile(ctx context.Context, userUUID string) (*models.User, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}

	user, err := s.repo.GetByUUID(ctx, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userUUID string, updates map[string]interface{}) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}

	user, err := s.repo.GetByUUID(ctx, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if firstName, ok := updates["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := updates["last_name"].(string); ok {
		user.LastName = lastName
	}
	if phoneNumber, ok := updates["phone_number"].(string); ok {
		user.PhoneNumber = phoneNumber
	}
	if username, ok := updates["username"].(string); ok {
		user.Username = username
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	log.Printf("UserService: Updated profile for user %s", userUUID)
	return nil
}

func (s *Service) GetCompanyDetails(ctx context.Context, userUUID string) (*models.CompanyDetails, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}

	company, err := s.repo.GetCompanyDetails(ctx, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company details not found")
		}
		return nil, fmt.Errorf("failed to get company details: %w", err)
	}

	return company, nil
}

func (s *Service) UpdateCompanyDetails(ctx context.Context, userUUID string, company *models.CompanyDetails) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}

	company.UserUUID = userUUID
	if err := s.repo.UpdateCompanyDetails(ctx, company); err != nil {
		return fmt.Errorf("failed to update company details: %w", err)
	}

	log.Printf("UserService: Updated company details for user %s", userUUID)
	return nil
}

func (s *Service) GetPaymentDetails(ctx context.Context, userUUID string) (*models.PaymentDetails, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}

	payment, err := s.repo.GetPaymentDetails(ctx, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment details not found")
		}
		return nil, fmt.Errorf("failed to get payment details: %w", err)
	}

	return payment, nil
}

func (s *Service) UpdatePaymentDetails(ctx context.Context, userUUID string, payment *models.PaymentDetails) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}

	payment.UserUUID = userUUID
	if err := s.repo.UpdatePaymentDetails(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment details: %w", err)
	}

	log.Printf("UserService: Updated payment details for user %s", userUUID)
	return nil
}

func (s *Service) GetPaymentHistory(ctx context.Context, userUUID string, limit, offset int) ([]models.PaymentHistory, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user UUID is required")
	}

	history, err := s.repo.GetPaymentHistory(ctx, userUUID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment history: %w", err)
	}

	return history, nil
}

func (s *Service) CreatePaymentHistory(ctx context.Context, userUUID string, payment *models.PaymentHistory) error {
	if userUUID == "" {
		return fmt.Errorf("user UUID is required")
	}

	payment.UserUUID = userUUID
	if err := s.repo.CreatePaymentHistory(ctx, payment); err != nil {
		return fmt.Errorf("failed to create payment history: %w", err)
	}

	log.Printf("UserService: Created payment history for user %s", userUUID)
	return nil
}

