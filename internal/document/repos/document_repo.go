package repos

import (
	"context"
	"time"

	"github.com/johnroshan2255/core-service/internal/document/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, doc *models.Document) error
	GetByID(ctx context.Context, id uint) (*models.Document, error)
	GetByUUID(ctx context.Context, userUUID string, id uint) (*models.Document, error)
	GetByUserUUID(ctx context.Context, userUUID string, limit, offset int) ([]models.Document, error)
	Update(ctx context.Context, doc *models.Document) error
	Delete(ctx context.Context, id uint) error

	GetExpiringDocuments(ctx context.Context, daysBeforeExpiry int) ([]models.Document, error)
	GetExpiredDocuments(ctx context.Context) ([]models.Document, error)
	UpdateNotificationSent(ctx context.Context, id uint, sent bool) error
}

type GORMRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{
		db: db,
	}
}

func (r *GORMRepository) Create(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *GORMRepository) GetByID(ctx context.Context, id uint) (*models.Document, error) {
	var doc models.Document
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *GORMRepository) GetByUUID(ctx context.Context, userUUID string, id uint) (*models.Document, error) {
	var doc models.Document
	if err := r.db.WithContext(ctx).Where("id = ? AND user_uuid = ?", id, userUUID).First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *GORMRepository) GetByUserUUID(ctx context.Context, userUUID string, limit, offset int) ([]models.Document, error) {
	var docs []models.Document
	query := r.db.WithContext(ctx).Where("user_uuid = ?", userUUID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GORMRepository) Update(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

func (r *GORMRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Document{}, id).Error
}

func (r *GORMRepository) GetExpiringDocuments(ctx context.Context, daysBeforeExpiry int) ([]models.Document, error) {
	var docs []models.Document
	thresholdDate := time.Now().AddDate(0, 0, daysBeforeExpiry)
	
	if err := r.db.WithContext(ctx).
		Where("expiry_date IS NOT NULL").
		Where("expiry_date <= ?", thresholdDate).
		Where("expiry_date > ?", time.Now()).
		Where("notification_sent = ?", false).
		Where("status = ?", models.DocumentStatusActive).
		Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GORMRepository) GetExpiredDocuments(ctx context.Context) ([]models.Document, error) {
	var docs []models.Document
	if err := r.db.WithContext(ctx).
		Where("expiry_date IS NOT NULL").
		Where("expiry_date <= ?", time.Now()).
		Where("status = ?", models.DocumentStatusActive).
		Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GORMRepository) UpdateNotificationSent(ctx context.Context, id uint, sent bool) error {
	return r.db.WithContext(ctx).Model(&models.Document{}).Where("id = ?", id).Update("notification_sent", sent).Error
}

