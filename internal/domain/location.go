package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(255);uniqueIndex:idx_name_user_id" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`

	ImageID int
	Image   Image

	UserID        int `gorm:"uniqueIndex:idx_name_user_id" json:"-"`
	LocationItems []LocationItem
}

type LocationRepository interface {
	Create(ctx context.Context, location Location) error
	Fetch(ctx context.Context, id int) (*Location, error)
	List(ctx context.Context, location Location) ([]*Location, error)
	Delete(ctx context.Context, location Location) error
}

type LocationUsecase interface {
	Create(ctx context.Context, location Location) error
	Fetch(ctx context.Context, id int) (*Location, error)
	List(ctx context.Context, location Location) ([]*Location, error)
	Delete(ctx context.Context, location Location) error
}
