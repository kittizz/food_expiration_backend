package usecase

import (
	"context"
	"sync"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type LocationUsecase struct {
	locationRepo domain.LocationRepository
	imageUsecase domain.ImageUsecase
}

func NewLocationUsecase(locationRepo domain.LocationRepository, imageUsecase domain.ImageUsecase) domain.LocationUsecase {
	return &LocationUsecase{
		locationRepo: locationRepo,
		imageUsecase: imageUsecase,
	}
}

// Create implements domain.LocationUsecase.
func (u *LocationUsecase) Create(ctx context.Context, location domain.Location) error {
	return u.locationRepo.Create(ctx, location)
}

func (u *LocationUsecase) Get(ctx context.Context, id int, items bool) (*domain.Location, error) {
	return u.locationRepo.Get(ctx, id, items)
}

func (u *LocationUsecase) List(ctx context.Context, userId int) ([]*domain.Location, error) {
	return u.locationRepo.List(ctx, userId)
}

func (u *LocationUsecase) Delete(ctx context.Context, location domain.Location) error {
	loc, err := u.locationRepo.Get(ctx, location.ID, true)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, item := range loc.Items {
		if loc.UserID == item.Image.UserID {
			wg.Add(1)
			go func(imageID int) {
				u.imageUsecase.Delete(ctx, imageID)
				wg.Done()
			}(item.ImageID)
		}
	}
	wg.Wait()

	err = u.locationRepo.Delete(ctx, location)
	if err != nil {
		return err
	}
	if loc.UserID == loc.Image.UserID {
		err := u.imageUsecase.Delete(ctx, loc.ImageID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *LocationUsecase) UpdateByID(ctx context.Context, location domain.Location, id int) error {
	loc, err := u.locationRepo.Get(ctx, id, false)
	if err != nil {
		return err
	}
	err = u.locationRepo.UpdateByID(ctx, location, id)
	if err != nil {
		return err
	}
	if loc.UserID == loc.Image.UserID && location.ImageID != loc.Image.ID {
		err := u.imageUsecase.Delete(ctx, loc.ImageID)
		if err != nil {
			return err
		}
	}

	return nil
}
