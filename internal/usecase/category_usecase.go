package usecase

import (
	"context"
	"strings"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type CategoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) domain.CategoryUsecase {
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

func (uc *CategoryUsecase) List(ctx context.Context) ([]string, error) {
	categorys, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, category := range categorys {
		names = append(names, category.Name)
	}
	return names, nil
}
func (uc *CategoryUsecase) Set(ctx context.Context, categories string) error {
	split := strings.Split(categories, "\\n")
	var categorys []*domain.Category
	for _, category := range split {
		categorys = append(categorys, &domain.Category{Name: category})
	}
	return uc.repo.Set(ctx, categorys)
}
