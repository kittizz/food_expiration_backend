package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type CategoryRepository struct {
	db *database.DatabaseMySQL
}

func NewCategoryRepository(db *database.DatabaseMySQL) domain.CategoryRepository {
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
func (repo *CategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	if err := repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) Set(ctx context.Context, categories []*domain.Category) error {
	err := repo.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Model(&domain.Category{}).
		Delete(&domain.Category{}).Error

	if err != nil {
		return err
	}
	err = repo.db.Create(&categories).Error
	return err

}
