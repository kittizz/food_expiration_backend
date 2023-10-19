package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type ItemUsecase struct {
	itemRepo domain.ItemRepository
}

func NewItemUsecase(itemRepo domain.ItemRepository) domain.ItemUsecase {
	return &ItemUsecase{
		itemRepo: itemRepo,
	}
}

// Create implements domain.ItemUsecase.
func (u *ItemUsecase) Create(ctx context.Context, item domain.Item) error {
	return u.itemRepo.Create(ctx, item)
}

func (u *ItemUsecase) Get(ctx context.Context, id int) (*domain.Item, error) {
	return u.itemRepo.Get(ctx, id)
}

func (u *ItemUsecase) Delete(ctx context.Context, item domain.Item) error {

	return u.itemRepo.Delete(ctx, item)
}

func (u *ItemUsecase) UpdateByID(ctx context.Context, item domain.Item, id int) error {

	return u.itemRepo.UpdateByID(ctx, item, id)
}

func (u *ItemUsecase) List(ctx context.Context, locationId *int, isArchived bool) ([]*domain.Item, error) {
	return u.itemRepo.List(ctx, locationId, isArchived)
}
