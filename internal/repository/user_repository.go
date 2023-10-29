package repository

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

type UserRepository struct {
	db *database.DatabaseMySQL
}

func NewUserRepository(db *database.DatabaseMySQL) domain.UserRepository {
	return &UserRepository{db}
}
func (repo *UserRepository) GetOrCreate(ctx context.Context, user domain.User) (*domain.User, error) {
	return &user, repo.db.WithContext(ctx).
		Where(domain.User{
			Uid: user.Uid,
		}).
		Assign(user).
		FirstOrCreate(&user).Error
}

// Get by uid,deviceid
func (repo *UserRepository) Get(ctx context.Context, user domain.User) (*domain.User, error) {

	err := repo.db.WithContext(ctx).
		Where(domain.User{
			Uid:      user.Uid,
			DeviceId: user.DeviceId,
		}).First(&user).Error
	return &user, err
}

// Update implements domain.UserRepository.
func (repo *UserRepository) UpdateByID(ctx context.Context, id int, user domain.User) error {
	return repo.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(user).Error
}
func (repo *UserRepository) ListNotifications(ctx context.Context, notiAt int) ([]*domain.User, error) {

	gap := notiAt - 10

	var result []*domain.User
	q := repo.db.WithContext(ctx).Model(domain.User{}).
		Select("id", "nickname", "fcm_token").
		Where("notification = ?", true).
		Where("fcm_token is not null")
	if gap < 0 {
		q = q.Where("notification_at <= ? or notification_at >= ?", notiAt, 1440+gap)
	} else {
		q = q.Where("notification_at <= ? or notification_at >= ?", notiAt, gap)
	}

	err := q.
		Find(&result).Error

	return result, err
}
