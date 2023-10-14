package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type ThumbnailUsecase struct {
	thumbnailRepo domain.ThumbnailRepository
}

func NewThumbnail(thumbnailRepo domain.ThumbnailRepository) domain.ThumbnailUsecase {
	return &ThumbnailUsecase{
		thumbnailRepo: thumbnailRepo,
	}
}

// Create implements domain.Thumbnail.
func (u *ThumbnailUsecase) Create(ctx context.Context, thum *domain.Thumbnail) error {
	return u.thumbnailRepo.Create(ctx, thum)
}

func (u *ThumbnailUsecase) GetByCategoryID(ctx context.Context, id int) (*domain.Thumbnail, error) {
	return u.thumbnailRepo.GetByCategoryID(ctx, id)
}

func (u *ThumbnailUsecase) ListByCategoryID(ctx context.Context, id int) ([]*domain.Thumbnail, error) {
	return u.thumbnailRepo.ListByCategoryID(ctx, id)
}

func (u *ThumbnailUsecase) Delete(ctx context.Context, thum domain.Thumbnail) error {
	return u.thumbnailRepo.Delete(ctx, thum)
}
