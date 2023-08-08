package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type BlogUsecase struct {
	blogRepo domain.BlogRepository
}

func NewBlogUsecase(blogRepo domain.BlogRepository) domain.BlogUsecase {
	return &BlogUsecase{
		blogRepo: blogRepo,
	}
}

func (u *BlogUsecase) List(ctx context.Context, isRandom bool, limit int) ([]*domain.Blog, error) {
	return u.blogRepo.List(ctx, isRandom, limit)
}
func (u *BlogUsecase) GetByID(ctx context.Context, id int) (*domain.Blog, error) {
	return u.blogRepo.GetByID(ctx, id)
}
