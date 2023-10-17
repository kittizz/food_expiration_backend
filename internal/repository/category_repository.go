package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type CategoryRepository struct {
	domain.CategoryRepository
	db *database.DatabaseMySQL
}

func NewCategoryRepository(db *database.DatabaseMySQL) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) Create(ctx context.Context, category domain.Category) error {
	return repo.db.Create(&category).Error
}

func (repo *CategoryRepository) Get(ctx context.Context, id int) (*domain.Category, error) {
	var category domain.Category
	if err := repo.db.
		Where(domain.Category{ID: id}).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (repo *CategoryRepository) Delete(ctx context.Context, category domain.Category) error {
	return repo.db.
		Unscoped().
		Delete(&category).Error
}

func (repo *CategoryRepository) Update(ctx context.Context, category domain.Category, id int) error {
	return repo.db.
		Where(domain.Category{ID: id}).
		Update("name", category.Name).Error
}
