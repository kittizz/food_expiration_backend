package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ThumbnailCategory struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name string `gorm:"type:varchar(255)" json:"name"`

	Type string `gorm:"type:varchar(8)" json:"type"`

	ImageID int   `json:"-"`
	Image   Image `json:"image"`

	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`
}

type ThumbnailCategoryRepository interface {
	Create(ctx context.Context, thumCategory *ThumbnailCategory) error
	Get(ctx context.Context, id int) (*ThumbnailCategory, error)
	List(ctx context.Context) ([]*ThumbnailCategory, error)
	Delete(ctx context.Context, thumCategory ThumbnailCategory) error
}

type ThumbnailCategoryUsecase interface {
	Create(ctx context.Context, thumCategory *ThumbnailCategory) error
	Get(ctx context.Context, id int) (*ThumbnailCategory, error)
	List(ctx context.Context) ([]*ThumbnailCategory, error)
	Delete(ctx context.Context, thumCategory ThumbnailCategory) error
}
