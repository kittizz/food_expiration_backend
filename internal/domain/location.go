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

	Name        string `gorm:"type:varchar(64);uniqueIndex:idx_name_user_id" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`

	ImageID int   `json:"-"`
	Image   Image `json:"image"`

	UserID    int        `gorm:"uniqueIndex:idx_name_user_id" json:"-"`
	Locations []Location `json:"locations,omitempty"`
}

type LocationRepository interface {
	Create(ctx context.Context, location Location) error
	Get(ctx context.Context, id int, items bool) (*Location, error)
	List(ctx context.Context, location Location) ([]*Location, error)
	Delete(ctx context.Context, location Location) error
	UpdateByID(ctx context.Context, location Location, id int) error
}

type LocationUsecase interface {
	Create(ctx context.Context, location Location) error
	Get(ctx context.Context, id int, items bool) (*Location, error)
	List(ctx context.Context, location Location) ([]*Location, error)
	Delete(ctx context.Context, location Location) error
	UpdateByID(ctx context.Context, location Location, id int) error
}
