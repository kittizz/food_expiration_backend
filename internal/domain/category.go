package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:varchar(32)" json:"name"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category Category) error
	Get(ctx context.Context, id int) (*Category, error)
	Delete(ctx context.Context, category Category) error
	Update(ctx context.Context, category Category, id int) error
	List(ctx context.Context) ([]*Category, error)
}

type CategoryUsecase interface {
	Create(ctx context.Context, category Category) error
	Get(ctx context.Context, id int) (*Category, error)
	Delete(ctx context.Context, category Category) error
	Update(ctx context.Context, category Category, id int) error
	List(ctx context.Context) ([]string, error)
}
