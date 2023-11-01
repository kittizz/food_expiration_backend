package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type AdminUsecase struct {
	locationRepo domain.LocationRepository
	userRepo     domain.UserRepository
	itemRepo     domain.ItemRepository
}

func NewAdminUsecase(locationRepo domain.LocationRepository, userRepo domain.UserRepository, itemRepo domain.ItemRepository) domain.AdminUsecase {
	return &AdminUsecase{
		locationRepo: locationRepo,
		userRepo:     userRepo,
		itemRepo:     itemRepo,
	}
}

func (u *AdminUsecase) Dashboard(ctx context.Context) (*domain.AdminDashboard, error) {
	locations, err := u.locationRepo.Counts(ctx)
	if err != nil {
		return nil, err
	}
	items, err := u.itemRepo.Counts(ctx)
	if err != nil {
		return nil, err
	}
	users, err := u.userRepo.Counts(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.AdminDashboard{
		Users:     users,
		Items:     items,
		Locations: locations,
	}, nil
}
