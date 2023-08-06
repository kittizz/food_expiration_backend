package domain

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        int            `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title   string `gorm:"type:varchar(255)"`
	Content string
	Image   string `gorm:"type:varchar(255)"`
}
