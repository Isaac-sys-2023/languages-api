package repository

import (
	"context"
	"errors"
	"languages-api/internal/models"

	"gorm.io/gorm"
)

type LanguageRepository interface {
	Create(ctx context.Context, language *models.Language) error
	CreateBatch(ctx context.Context, languages []models.Language) error
	FindByID(ctx context.Context, id uint) (*models.Language, error)
	FindAll(ctx context.Context, page, pageSize int) ([]models.Language, int64, error)
	FindWithFilters(ctx context.Context, filters LanguageFilters, page, pageSize int) ([]models.Language, int64, error)
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
	Save(ctx context.Context, user *models.Language) error
	Delete(ctx context.Context, id uint) error
	HardDelete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
}

type languageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) LanguageRepository {
	return &languageRepository{db: db}
}

func (r *languageRepository) Create(ctx context.Context, language *models.Language) error {
	result := r.db.WithContext(ctx).Create(language)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *languageRepository) CreateBatch(ctx context.Context, languages []models.Language) error {
	result := r.db.WithContext(ctx).CreateInBatches(languages, 100)
	return result.Error
}

func (r *languageRepository) FindByID(ctx context.Context, id uint) (*models.Language, error) {
	var language models.Language

	result := r.db.WithContext(ctx).First(&language, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found
		}
		return nil, result.Error
	}
	return &language, nil
}

func (r *languageRepository) FindAll(ctx context.Context, page, pageSize int) ([]models.Language, int64, error) {
	var languages []models.Language
	var total int64

	r.db.WithContext(ctx).Model(&models.Language{}).Count(&total)

	offset := (page - 1) * pageSize

	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&languages)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return languages, total, nil
}

func (r *languageRepository) FindWithFilters(ctx context.Context, filters LanguageFilters, page, pageSize int) ([]models.Language, int64, error) {
	var languages []models.Language
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Language{})

	if filters.Creator != "" {
		query = query.Where("creator = ?", filters.Creator)
	}
	if filters.Name != "" {
		query = query.Where("name = ?", filters.Name)
	}
	if filters.ReleaseYear > 0 {
		query = query.Where("release_year > ?", filters.ReleaseYear)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&languages).
		Error

	if err != nil {
		return nil, 0, err
	}

	return languages, total, nil
}

// LanguageFilters contains optional filters for querying users
type LanguageFilters struct {
	Creator     string
	Name        string
	ReleaseYear int
}

func (r *languageRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	// Updates only the specified fields
	result := r.db.WithContext(ctx).
		Model(&models.Language{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *languageRepository) Save(ctx context.Context, user *models.Language) error {
	// Save will update all fields, including zero values
	// Use this when you want to explicitly set fields to zero/empty
	result := r.db.WithContext(ctx).Save(user)
	return result.Error
}

// Delete performs a soft delete (sets deleted_at)
func (r *languageRepository) Delete(ctx context.Context, id uint) error {
	// With gorm.Model, Delete sets deleted_at instead of removing the row
	result := r.db.WithContext(ctx).Delete(&models.Language{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// HardDelete permanently removes a record from the database
func (r *languageRepository) HardDelete(ctx context.Context, id uint) error {
	// Unscoped bypasses soft delete and permanently removes the record
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.Language{}, id)
	return result.Error
}

// Restore recovers a soft-deleted record
func (r *languageRepository) Restore(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Model(&models.Language{}).
		Where("id = ?", id).
		Update("deleted_at", nil)
	return result.Error
}
