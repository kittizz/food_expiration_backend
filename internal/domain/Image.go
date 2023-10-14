package domain

import (
	"context"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Path     string `gorm:"type:varchar(255)" json:"path"`
	BlurHash string `gorm:"type:varchar(30)" json:"blurHash"`

	UserID int  `json:"-"`
	User   User `json:"-"`
}

type ImageRepository interface {
	Create(ctx context.Context, image *Image) error
	Get(ctx context.Context, id int) (*Image, error)
	Delete(ctx context.Context, id int) error
}

type ImageUsecase interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, hash, folder string, userId int) (img *Image, err error)
	Delete(ctx context.Context, id int) error
	DeleteWithPath(ctx context.Context, path string) error
}
