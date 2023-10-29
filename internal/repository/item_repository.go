package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type ItemRepository struct {
	db *database.DatabaseMySQL
}

func NewItemRepository(db *database.DatabaseMySQL) domain.ItemRepository {
	return &ItemRepository{db}
}
func (repo *ItemRepository) Create(ctx context.Context, item domain.Item) error {
	return repo.db.WithContext(ctx).
		Create(&item).Error
}

func (repo *ItemRepository) Get(ctx context.Context, id int) (*domain.Item, error) {
	var result domain.Item
	q := repo.db.WithContext(ctx).
		Where(domain.Item{ID: id})

	err := q.
		Joins("Image").
		First(&result).Error

	return &result, err
}

func (repo *ItemRepository) Delete(ctx context.Context, item domain.Item) error {
	return repo.db.WithContext(ctx).Unscoped().
		Delete(&item).Error
}
func (repo *ItemRepository) UpdateByID(ctx context.Context, item domain.Item, id int) error {
	return repo.db.WithContext(ctx).
		Where(domain.Item{ID: id}).
		Updates(item).Error
}

func (repo *ItemRepository) List(ctx context.Context, userId int, locationId int, isArchived bool) ([]*domain.Item, error) {
	var result []*domain.Item
	q := repo.db.WithContext(ctx).Model(domain.Item{}).
		Where("is_archived = ?", isArchived).
		Where(domain.Item{UserID: userId, LocationID: locationId})

	err := q.
		Joins("Image").
		Find(&result).Error

	return result, err
}
func (repo *ItemRepository) Archive(ctx context.Context, archive bool, id []int) error {
	return repo.db.WithContext(ctx).
		Model(&domain.Item{}).
		Where("id IN ?", id).
		Update("is_archived", archive).Error
}

func (repo *ItemRepository) ListForNotification(ctx context.Context, users []int) ([]*domain.Item, error) {
	var result []*domain.Item
	q := repo.db.WithContext(ctx).Model(domain.Item{}).
		Where("is_archived = ?", false).
		Where("user_id IN (?)", users).
		Where("notification_status <> ?", domain.NOTIFICATION_STATUS_DONE).
		Where("DATE_ADD(expire_date,INTERVAL -forewarn_day DAY) <= CURDATE()").
		Where("last_notification_at < CURDATE() or last_notification_at is null")

	err := q.
		Select("id", "name", "user_id", "forewarn_day", "expire_date", "image_id", "notification_status").
		Find(&result).Error

	return result, err
}
func (repo *ItemRepository) UpdateNotificationStatus(ctx context.Context, id []int, status domain.NotificationStatus) error {
	return repo.db.WithContext(ctx).
		Model(&domain.Item{}).
		Where("id IN (?)", id).
		Updates(map[string]any{
			"notification_status":  status,
			"last_notification_at": gorm.Expr("CURDATE()"),
		}).
		Error
}

func (repo *ItemRepository) Deletes(ctx context.Context, ids []int) error {
	return repo.db.WithContext(ctx).Unscoped().Delete(&domain.Item{}, ids).Error
}
