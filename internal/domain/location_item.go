package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type LocationItem struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`

	StorageDate time.Time `json:"storageDate"`
	ExpireDate  time.Time `json:"expireDate"`

	ForewarnDay int    `gorm:"type:int" json:"forewarnDay"`
	IsArchived  bool   `gorm:"type:boolean" json:"isArchived"`
	Category    string `gorm:"type:varchar(255)" json:"category"`
	Barcode     string `gorm:"type:varchar(255)" json:"barcode"`

	LocationID int   `json:"-"`
	ImageID    int   `json:"-"`
	Image      Image `json:"image"`
}

type LocationItemRepository interface {
	Create(ctx context.Context, locationItem LocationItem) error
	Get(ctx context.Context, id int) (*LocationItem, error)
	Delete(ctx context.Context, locationItem LocationItem) error
	UpdateByID(ctx context.Context, locationItem LocationItem, id int) error
}

type LocationItemUsecase interface {
	Create(ctx context.Context, locationItem LocationItem) error
	Get(ctx context.Context, id int) (*LocationItem, error)
	Delete(ctx context.Context, locationItem LocationItem) error
	UpdateByID(ctx context.Context, locationItem LocationItem, id int) error
}
