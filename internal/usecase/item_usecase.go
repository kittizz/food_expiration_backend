package usecase

import (
	"context"
	"sort"

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

func (u *ItemUsecase) List(ctx context.Context, userId int, locationId int, isArchived bool, isSort bool) ([]*domain.Item, error) {

	items, err := u.itemRepo.List(ctx, userId, locationId, isArchived)
	if err != nil {
		return nil, err
	}
	if isSort {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].ExpireDate.AddDate(0, 0, -*items[i].ForewarnDay).Before(items[j].ExpireDate.AddDate(0, 0, -*items[j].ForewarnDay))
		})
	}
	return items, nil

}

func (u *ItemUsecase) Archive(ctx context.Context, archive bool, id []int) error {

	return u.itemRepo.Archive(ctx, archive, id)
}
func (u *ItemUsecase) ListForNotification(ctx context.Context, users []int) ([]*domain.Item, error) {
	return u.itemRepo.ListForNotification(ctx, users)
}
func (u *ItemUsecase) UpdateNotificationStatus(ctx context.Context, notiMap map[domain.NotificationStatus][]int) error {
	for toStatus, ids := range notiMap {
		err := u.itemRepo.UpdateNotificationStatus(ctx, ids, toStatus)
		if err != nil {
			return err
		}
	}
	return nil
}
