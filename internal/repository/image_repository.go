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
		Create(&image).Error
}
