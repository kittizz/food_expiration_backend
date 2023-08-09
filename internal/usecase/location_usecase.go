package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type LocationUsecase struct {
	locationRepo domain.LocationRepository
}

func NewLocationUsecase(locationRepo domain.LocationRepository) domain.LocationUsecase {
	return &LocationUsecase{
		locationRepo: locationRepo,
	}
}

// Create implements domain.LocationUsecase.
func (u *LocationUsecase) Create(ctx context.Context, location domain.Location) error {
	return u.locationRepo.Create(ctx, location)
}

func (u *LocationUsecase) Fetch(ctx context.Context, id int) (*domain.Location, error) {
	return u.locationRepo.Fetch(ctx, id)
}

func (u *LocationUsecase) List(ctx context.Context, location domain.Location) ([]*domain.Location, error) {
	return u.locationRepo.List(ctx, location)
}

func (u *LocationUsecase) Delete(ctx context.Context, location domain.Location) error {
	return u.locationRepo.Delete(ctx, location)
}
