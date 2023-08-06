package domain

import (
	"time"

	"gorm.io/gorm"
)

type ThumbnailCategory struct {
	ID        int            `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name  string `gorm:"type:varchar(255)"`
	Image string `gorm:"type:varchar(255)"`

	Thumbnails []Thumbnail
}
