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
func (r *BlogRepository) Update(ctx context.Context, blog domain.Blog, id int) (int, error) {
	q := r.db.WithContext(ctx).Model(&blog).Where(domain.Blog{
		ID: id,
	})

	if id != 0 {
		err := q.Updates(domain.Blog{
			Title:   blog.Title,
			Content: blog.Content,
			ImageID: blog.ImageID,
		}).Error
		if err != nil {
			return blog.ID, err
		}
		return blog.ID, nil
	}
	err := q.Create(&blog).Error
	if err != nil {
		return blog.ID, err
	}
	return blog.ID, nil

}
func (r *BlogRepository) Delete(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Unscoped().Where(domain.Blog{
		ID: id,
	}).Delete(&domain.Blog{}).Error
	if err != nil {
		return err
	}
	return nil
}
