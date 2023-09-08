package domain

import (
	"context"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        int            `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Path     string `gorm:"type:varchar(255)"`
	BlurHash string `gorm:"type:varchar(30)"`

	UserID int  `json:"-"`
	User   User `json:"-"`
}

type ImageRepository interface {
	Create(ctx context.Context, image Image) error
}

type ImageUsecase interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, userId int) error
}
