package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type ThumbnailCategoryRepository struct {
	db *database.DatabaseMySQL
}

func NewThumbnailCategoryRepository(db *database.DatabaseMySQL) domain.ThumbnailCategoryRepository {
	return &ThumbnailCategoryRepository{db}
}
func (repo *ThumbnailCategoryRepository) Create(ctx context.Context, thumCategory *domain.ThumbnailCategory) error {
	return repo.db.WithContext(ctx).
		Create(thumCategory).Error
}

func (repo *ThumbnailCategoryRepository) Get(ctx context.Context, id int) (*domain.ThumbnailCategory, error) {
	var result domain.ThumbnailCategory
	err := repo.db.WithContext(ctx).
		Where(&domain.ThumbnailCategory{
			ID: id,
		}).
		Preload("Thumbnails").
		Preload("Thumbnails.Image").
		Joins("Image").
		First(&result).Error
	return &result, err
}

func (repo *ThumbnailCategoryRepository) List(ctx context.Context) ([]*domain.ThumbnailCategory, error) {
	var result []*domain.ThumbnailCategory
	err := repo.db.WithContext(ctx).
		Joins("Image").
		Find(&result).Error
	return result, err
}
func (repo *ThumbnailCategoryRepository) Delete(ctx context.Context, thumCategory domain.ThumbnailCategory) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&thumCategory).Error
}
func (repo *ThumbnailCategoryRepository) Update(ctx context.Context, thumCategory domain.ThumbnailCategory, id int) error {
	return repo.db.WithContext(ctx).
		Model(&thumCategory).
		Where(domain.ThumbnailCategory{
			ID: id,
		}).
		Updates(&thumCategory).Error
}
