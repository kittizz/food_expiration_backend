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

func (repo *LocationRepository) Fetch(ctx context.Context, id int) (*domain.Location, error) {
	var result domain.Location
	err := repo.db.WithContext(ctx).
		Where(domain.Location{ID: id}).
		First(&result).Error
	return &result, err
}

func (repo *LocationRepository) List(ctx context.Context, query domain.Location) ([]*domain.Location, error) {
	var result []*domain.Location
	err := repo.db.WithContext(ctx).
		Where(domain.Location{ID: query.ID, UserID: query.UserID}).
		First(&result).Error
	return result, err
}
func (repo *LocationRepository) Delete(ctx context.Context, location domain.Location) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&location).Error
}
