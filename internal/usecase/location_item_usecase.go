package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type LocationItemUsecase struct {
	locationItemRepo domain.LocationItemRepository
}

func NewLocationItemUsecase(locationItemRepo domain.LocationItemRepository) domain.LocationItemUsecase {
	return &LocationItemUsecase{
		locationItemRepo: locationItemRepo,
	}
}

// Create implements domain.LocationItemUsecase.
func (u *LocationItemUsecase) Create(ctx context.Context, locationItem domain.LocationItem) error {
	return u.locationItemRepo.Create(ctx, locationItem)
}

func (u *LocationItemUsecase) Get(ctx context.Context, id int) (*domain.LocationItem, error) {
	return u.locationItemRepo.Get(ctx, id)
}

func (u *LocationItemUsecase) Delete(ctx context.Context, locationItem domain.LocationItem) error {

	return u.locationItemRepo.Delete(ctx, locationItem)
}

func (u *LocationItemUsecase) UpdateByID(ctx context.Context, locationItem domain.LocationItem, id int) error {

	return u.locationItemRepo.UpdateByID(ctx, locationItem, id)
}
