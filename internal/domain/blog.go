package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title   string `gorm:"type:varchar(255)" json:"title"`
	Content string `gorm:"type:text" json:"content"`

	ImageID int   `json:"-"`
	Image   Image `json:"image"`
}

type BlogRepository interface {
	List(ctx context.Context, isRandom bool, limit int) ([]*Blog, error)
	GetByID(ctx context.Context, id int) (*Blog, error)
	Update(ctx context.Context, blog Blog, id int) (int, error)
	Delete(ctx context.Context, id int) error
}

type BlogUsecase interface {
	List(ctx context.Context, isRandom bool, limit int) ([]*Blog, error)
	GetByID(ctx context.Context, id int) (*Blog, error)
	Rename(ctx context.Context, name string, id int) error
	UpdateOrCreate(ctx context.Context, blog Blog, id int) (int, error)
	Delete(ctx context.Context, id int) error
}
