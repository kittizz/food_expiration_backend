package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type BlogRepository struct {
	db *database.DatabaseMySQL
}

func NewBlogRepository(db *database.DatabaseMySQL) domain.BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) List(ctx context.Context, isRandom bool, limit int) ([]*domain.Blog, error) {
	var blogs []*domain.Blog

	q := r.db.WithContext(ctx)
	if isRandom {
		q = q.Order("RAND()")
	} else {
		q = q.Order("id DESC")
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Joins("Image").Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *BlogRepository) GetByID(ctx context.Context, id int) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.db.WithContext(ctx).Where(domain.Blog{
		ID: id,
	}).Joins("Image").First(&blog).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}
