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
func (repo *UserRepository) FetchOrCreate(ctx context.Context, user domain.User) (*domain.User, error) {
	return &user, repo.db.WithContext(ctx).
		Where(domain.User{
			Uid: user.Uid,
		}).
		Assign(user).
		FirstOrCreate(&user).Error
}
