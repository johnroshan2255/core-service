package repos

import (
	"context"

	"github.com/johnroshan2255/core-service/internal/user/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByUUID(ctx context.Context, userUUID string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userUUID string) error

	GetCompanyDetails(ctx context.Context, userUUID string) (*models.CompanyDetails, error)
	UpdateCompanyDetails(ctx context.Context, company *models.CompanyDetails) error

	GetPaymentDetails(ctx context.Context, userUUID string) (*models.PaymentDetails, error)
	UpdatePaymentDetails(ctx context.Context, payment *models.PaymentDetails) error

	GetPaymentHistory(ctx context.Context, userUUID string, limit, offset int) ([]models.PaymentHistory, error)
	CreatePaymentHistory(ctx context.Context, payment *models.PaymentHistory) error
}

type GORMRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{
		db: db,
	}
}

func (r *GORMRepository) GetByUUID(ctx context.Context, userUUID string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GORMRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GORMRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GORMRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Model(user).Where("uuid = ?", user.UUID).Updates(user).Error
}

func (r *GORMRepository) Delete(ctx context.Context, userUUID string) error {
	return r.db.WithContext(ctx).Where("uuid = ?", userUUID).Delete(&models.User{}).Error
}

func (r *GORMRepository) GetCompanyDetails(ctx context.Context, userUUID string) (*models.CompanyDetails, error) {
	var company models.CompanyDetails
	if err := r.db.WithContext(ctx).Where("user_uuid = ?", userUUID).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *GORMRepository) UpdateCompanyDetails(ctx context.Context, company *models.CompanyDetails) error {
	var existing models.CompanyDetails
	err := r.db.WithContext(ctx).Where("user_uuid = ?", company.UserUUID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(company).Error
	}
	if err != nil {
		return err
	}
	company.ID = existing.ID
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *GORMRepository) GetPaymentDetails(ctx context.Context, userUUID string) (*models.PaymentDetails, error) {
	var payment models.PaymentDetails
	if err := r.db.WithContext(ctx).Where("user_uuid = ?", userUUID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *GORMRepository) UpdatePaymentDetails(ctx context.Context, payment *models.PaymentDetails) error {
	var existing models.PaymentDetails
	err := r.db.WithContext(ctx).Where("user_uuid = ?", payment.UserUUID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(payment).Error
	}
	if err != nil {
		return err
	}
	payment.ID = existing.ID
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *GORMRepository) GetPaymentHistory(ctx context.Context, userUUID string, limit, offset int) ([]models.PaymentHistory, error) {
	var history []models.PaymentHistory
	query := r.db.WithContext(ctx).Where("user_uuid = ?", userUUID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *GORMRepository) CreatePaymentHistory(ctx context.Context, payment *models.PaymentHistory) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

