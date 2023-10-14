package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type ImageRepository struct {
	db *database.DatabaseMySQL
}

func NewImageRepository(db *database.DatabaseMySQL) domain.ImageRepository {
	return &ImageRepository{db: db}
}

// Create implements domain.ImageRepository.
func (repo *ImageRepository) Create(ctx context.Context, image *domain.Image) error {
	return repo.db.WithContext(ctx).
		Create(image).Error
}
func (repo *ImageRepository) Delete(ctx context.Context, id int) error {
	return repo.db.WithContext(ctx).
		Unscoped().
		Delete(&domain.Image{}, id).Error
}
func (repo *ImageRepository) Get(ctx context.Context, id int) (*domain.Image, error) {
	var result domain.Image
	err := repo.db.WithContext(ctx).
		Where(domain.Location{ID: id}).
		First(&result).Error
	return &result, err
}
