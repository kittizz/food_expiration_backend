package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type ThumbnailRepository struct {
	db *database.DatabaseMySQL
}

func NewThumbnailRepository(db *database.DatabaseMySQL) domain.ThumbnailRepository {
	return &ThumbnailRepository{db}
}
func (repo *ThumbnailRepository) Create(ctx context.Context, thum *domain.Thumbnail) error {
	return repo.db.WithContext(ctx).
		Create(thum).Error
}

func (repo *ThumbnailRepository) GetByCategoryID(ctx context.Context, id int) (*domain.Thumbnail, error) {
	var result domain.Thumbnail
	err := repo.db.WithContext(ctx).
		Where(domain.Thumbnail{ThumbnailCategoryID: id}).
		Joins("Image").
		First(&result).Error
	return &result, err
}

func (repo *ThumbnailRepository) ListByCategoryID(ctx context.Context, id int) ([]*domain.Thumbnail, error) {
	var result []*domain.Thumbnail
	err := repo.db.WithContext(ctx).
		Where(domain.Thumbnail{ThumbnailCategoryID: id}).
		Joins("Image").
		Find(&result).Error
	return result, err
}
func (repo *ThumbnailRepository) Delete(ctx context.Context, thum domain.Thumbnail) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&thum).Error
}
