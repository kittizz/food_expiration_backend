package domain

import (
	"context"
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

type NotificationStatus string

const (
	NOTIFICATION_STATUS_PLANNED  NotificationStatus = "PLANNED"
	NOTIFICATION_STATUS_FOREWARN NotificationStatus = "FOREWARN"
	NOTIFICATION_STATUS_EXPIRING NotificationStatus = "EXPIRING"
	NOTIFICATION_STATUS_DONE     NotificationStatus = "DONE"
)

func (n *NotificationStatus) Scan(value any) error {
	*n = NotificationStatus(value.([]byte))
	return nil
}

func (n NotificationStatus) Value() (driver.Value, error) {
	return string(n), nil
}

type Item struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        *string `gorm:"type:varchar(255)" json:"name"`
	Description *string `gorm:"type:varchar(255)" json:"description"`

	StorageDate time.Time `json:"storageDate"`
	ExpireDate  time.Time `json:"expireDate"`

	ForewarnDay        *int               `gorm:"type:int" json:"forewarnDay"`
	IsArchived         *bool              `json:"isArchived"`
	Category           *string            `gorm:"type:varchar(255)" json:"category"`
	Barcode            *string            `gorm:"type:varchar(255)" json:"barcode"`
	Quantity           *int               ` json:"quantity"`
	Unit               *string            `json:"unit"`
	NotificationStatus NotificationStatus `gorm:"type:enum('PLANNED', 'FOREWARN','EXPIRING', 'DONE');default:'PLANNED';index" json:"notificationStatus"`
	LastNotificationAt *time.Time         `json:"lastNotificationAt"`
	LocationID         int                `json:"locationId"`
	ImageID            int                `json:"-"`
	Image              Image              `json:"image"`

	UserID int  `json:"userId"`
	User   User `json:"-"`
}

type ItemRepository interface {
	Create(ctx context.Context, item Item) error
	Get(ctx context.Context, id int) (*Item, error)
	List(ctx context.Context, userId int, locationId int, isArchived bool) ([]*Item, error)
	ListForNotification(ctx context.Context, users []int) ([]*Item, error)
	Delete(ctx context.Context, item Item) error
	Deletes(ctx context.Context, ids []int) error
	UpdateByID(ctx context.Context, item Item, id int) error
	UpdateNotificationStatus(ctx context.Context, id []int, status NotificationStatus) error
	Archive(ctx context.Context, archive bool, id []int) error
	Counts(ctx context.Context) (int64, error)
}

type ItemUsecase interface {
	Create(ctx context.Context, item Item) error
	Get(ctx context.Context, id int) (*Item, error)
	List(ctx context.Context, userId int, locationId int, isArchived bool, sort bool) ([]*Item, error)
	ListForNotification(ctx context.Context, users []int) ([]*Item, error)
	Delete(ctx context.Context, item Item) error
	Deletes(ctx context.Context, ids []int) error
	UpdateByID(ctx context.Context, item Item, id int) error
	Archive(ctx context.Context, archive bool, id []int) error
	UpdateNotificationStatus(ctx context.Context, notiMap map[NotificationStatus][]int) error
}
