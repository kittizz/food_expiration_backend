package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type ThumbnailCategoryUsecase struct {
	thumbnailCategoryRepo domain.ThumbnailCategoryRepository
}

func NewThumbnailCategoryUsecase(thumbnailCategoryRepo domain.ThumbnailCategoryRepository) domain.ThumbnailCategoryUsecase {
	return &ThumbnailCategoryUsecase{
		thumbnailCategoryRepo: thumbnailCategoryRepo,
	}
}

// Create implements domain.ThumbnailCategory.
func (u *ThumbnailCategoryUsecase) Create(ctx context.Context, thumCategory *domain.ThumbnailCategory) error {
	return u.thumbnailCategoryRepo.Create(ctx, thumCategory)
}

func (u *ThumbnailCategoryUsecase) Get(ctx context.Context, id int) (*domain.ThumbnailCategory, error) {
	return u.thumbnailCategoryRepo.Get(ctx, id)
}

func (u *ThumbnailCategoryUsecase) List(ctx context.Context) ([]*domain.ThumbnailCategory, error) {
	return u.thumbnailCategoryRepo.List(ctx)
}

func (u *ThumbnailCategoryUsecase) Delete(ctx context.Context, thumCategory domain.ThumbnailCategory) error {
	return u.thumbnailCategoryRepo.Delete(ctx, thumCategory)
}
func (u *ThumbnailCategoryUsecase) Update(ctx context.Context, thumCategory domain.ThumbnailCategory, id int) error {
	return u.thumbnailCategoryRepo.Update(ctx, thumCategory, id)
}
