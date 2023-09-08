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
	Uid            string `gorm:"unique;type:varchar(64)" json:"-"`

	Role string `gorm:"type:varchar(64);default:'user'" json:"role"`

	DeviceId *string `gorm:"type:varchar(36);uniqueIndex;" json:"-"`

	Locations []Location `json:"locations,omitempty" `
}

type UserUsecase interface {
	RegisterDevice(ctx context.Context, id int) (string, error)
	GetAuthUserByUid(ctx context.Context, uid string) (*User, error)
	VerifyIDToken(ctx context.Context, authorization string) (*User, error)
	GenerateIDToken(ctx context.Context, uid string) (string, error)
	Sync(ctx context.Context, user User) (*User, error)
	GetUserByDeviceId(ctx context.Context, deviceId string) (*User, error)
}

type UserRepository interface {
	FetchOrCreate(ctx context.Context, user User) (*User, error)
	Fetch(ctx context.Context, user User) (*User, error)
	UpdateByID(ctx context.Context, id int, user User) error
}
