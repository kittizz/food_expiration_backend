package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Item struct {
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

type ItemRepository interface {
	Create(ctx context.Context, item Item) error
	Get(ctx context.Context, id int) (*Item, error)
	Delete(ctx context.Context, item Item) error
	UpdateByID(ctx context.Context, item Item, id int) error
}

type ItemUsecase interface {
	Create(ctx context.Context, item Item) error
	Get(ctx context.Context, id int) (*Item, error)
	Delete(ctx context.Context, item Item) error
	UpdateByID(ctx context.Context, item Item, id int) error
}