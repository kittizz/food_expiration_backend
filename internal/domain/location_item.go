package domain

import (
	"time"

	"gorm.io/gorm"
)

type LocationItem struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`

	ImageID int
	Image   Image

	ExpiryDate time.Time

	IsArchived  bool `gorm:"type:boolean"`
	ForewarnDay int  `gorm:"type:int"`

	LocationID int
}
