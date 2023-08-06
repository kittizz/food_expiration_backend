package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID        int            `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(255);uniqueIndex:idx_name_user_id"`
	Description string `gorm:"type:varchar(255)"`
	Image       string `gorm:"type:varchar(255)"`

	UserID        int `gorm:"uniqueIndex:idx_name_user_id"`
	LocationItems []LocationItem
}

type LocationRepository interface {
	Create(ctx context.Context, location *Location) error
	Fetch(ctx context.Context, id int) (*Location, error)
}

type LocationUsecase interface {
	Create(ctx context.Context, location *Location) error
}
