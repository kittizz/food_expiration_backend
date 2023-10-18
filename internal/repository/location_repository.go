package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type LocationRepository struct {
	db *database.DatabaseMySQL
}

func NewLocationRepository(db *database.DatabaseMySQL) domain.LocationRepository {
	return &LocationRepository{db}
}
func (repo *LocationRepository) Create(ctx context.Context, location domain.Location) error {
	return repo.db.WithContext(ctx).
		Create(&location).Error
}

func (repo *LocationRepository) Get(ctx context.Context, id int, items bool) (*domain.Location, error) {
	var result domain.Location
	q := repo.db.WithContext(ctx).
		Where(domain.Location{ID: id})

	if items {
		q.
			Preload("Items").
			Preload("Items.Image")
	}
	err := q.
		Joins("Image").
		First(&result).Error

	return &result, err
}

func (repo *LocationRepository) List(ctx context.Context, query domain.Location) ([]*domain.Location, error) {
	var result []*domain.Location
	err := repo.db.WithContext(ctx).
		Where(domain.Location{ID: query.ID, UserID: query.UserID}).
		Order("id DESC").
		Joins("Image").
		Find(&result).Error
	return result, err
}
func (repo *LocationRepository) Delete(ctx context.Context, location domain.Location) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&location).Error
}
func (repo *LocationRepository) UpdateByID(ctx context.Context, location domain.Location, id int) error {
	return repo.db.WithContext(ctx).
		Where(domain.Location{ID: id}).
		Updates(location).Error
}
