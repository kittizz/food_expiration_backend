package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Thumbnail struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name string `gorm:"type:varchar(255)" json:"name"`

	ImageID int   `json:"-"`
	Image   Image `json:"image"`

	ThumbnailCategoryID int `json:"thumbnailCategoryId"`
}

type ThumbnailRepository interface {
	Create(ctx context.Context, thum *Thumbnail) error
	GetByCategoryID(ctx context.Context, id int) (*Thumbnail, error)
	ListByCategoryID(ctx context.Context, id int) ([]*Thumbnail, error)
	Delete(ctx context.Context, thum Thumbnail) error
}

type ThumbnailUsecase interface {
	Create(ctx context.Context, thum *Thumbnail) error
	GetByCategoryID(ctx context.Context, id int) (*Thumbnail, error)
	ListByCategoryID(ctx context.Context, id int) ([]*Thumbnail, error)
	Delete(ctx context.Context, thum Thumbnail) error
}
