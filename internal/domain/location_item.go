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

	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`

	ImageID int   `json:"-"`
	Image   Image `json:"image"`

	ExpiryDate time.Time `json:"expiryDate"`

	IsArchived  bool `gorm:"type:boolean" json:"isArchived"`
	ForewarnDay int  `gorm:"type:int" json:"forewarnDay"`

	LocationID int `json:"locationID"`
}
