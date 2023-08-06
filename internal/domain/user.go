package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int            `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email          string `gorm:"type:varchar(255)" json:"email"`
	SignInProvider string `gorm:"type:varchar(64)" json:"signInProvider"`
	Uid            string `gorm:"unique;type:varchar(64)" json:"uid"`

	Locations []Location `json:"locations,omitempty" `
}

type UserUsecase interface {
	GetAuthUserByUid(ctx context.Context, uid string) (*User, error)
	VerifyIDToken(ctx context.Context, authorization string) (*User, error)
	GenerateIDToken(ctx context.Context, uid string) (string, error)
	Sync(ctx context.Context, user User) (*User, error)
}

type UserRepository interface {
	FetchOrCreate(ctx context.Context, user User) (*User, error)
}
