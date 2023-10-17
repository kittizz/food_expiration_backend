package usecase

import (
	"context"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type CategoryUsecase struct {
	domain.CategoryUsecase
	repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
	}
}

func (uc *CategoryUsecase) Create(ctx context.Context, category domain.Category) error {
	return uc.repo.Create(ctx, category)
}

func (uc *CategoryUsecase) Get(ctx context.Context, id int) (*domain.Category, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *CategoryUsecase) Delete(ctx context.Context, category domain.Category) error {
	return uc.repo.Delete(ctx, category)
}

func (uc *CategoryUsecase) Update(ctx context.Context, category domain.Category, id int) error {
	return uc.repo.Update(ctx, category, id)
}
