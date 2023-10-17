package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type LocationItemRepository struct {
	db *database.DatabaseMySQL
}

func NewLocationItemRepository(db *database.DatabaseMySQL) domain.LocationItemRepository {
	return &LocationItemRepository{db}
}
func (repo *LocationItemRepository) Create(ctx context.Context, locationItem domain.LocationItem) error {
	return repo.db.WithContext(ctx).
		Create(&locationItem).Error
}

func (repo *LocationItemRepository) Get(ctx context.Context, id int) (*domain.LocationItem, error) {
	var result domain.LocationItem
	q := repo.db.WithContext(ctx).
		Where(domain.LocationItem{ID: id})

	err := q.
		Joins("Image").
		First(&result).Error

	return &result, err
}

func (repo *LocationItemRepository) Delete(ctx context.Context, locationItem domain.LocationItem) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&locationItem).Error
}
func (repo *LocationItemRepository) UpdateByID(ctx context.Context, locationItem domain.LocationItem, id int) error {
	return repo.db.WithContext(ctx).
		Where(domain.LocationItem{ID: id}).
		Updates(locationItem).Error
}
