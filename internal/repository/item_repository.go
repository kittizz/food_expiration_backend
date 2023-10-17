package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type ItemRepository struct {
	db *database.DatabaseMySQL
}

func NewItemRepository(db *database.DatabaseMySQL) domain.ItemRepository {
	return &ItemRepository{db}
}
func (repo *ItemRepository) Create(ctx context.Context, item domain.Item) error {
	return repo.db.WithContext(ctx).
		Create(&item).Error
}

func (repo *ItemRepository) Get(ctx context.Context, id int) (*domain.Item, error) {
	var result domain.Item
	q := repo.db.WithContext(ctx).
		Where(domain.Item{ID: id})

	err := q.
		Joins("Image").
		First(&result).Error

	return &result, err
}

func (repo *ItemRepository) Delete(ctx context.Context, item domain.Item) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&item).Error
}
func (repo *ItemRepository) UpdateByID(ctx context.Context, item domain.Item, id int) error {
	return repo.db.WithContext(ctx).
		Where(domain.Item{ID: id}).
		Updates(item).Error
}
