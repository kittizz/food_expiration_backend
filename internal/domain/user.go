package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type NotificationAt struct {
	time.Duration
	json.Marshaler
	json.Unmarshaler
}

func (n NotificationAt) MarshalJSON() ([]byte, error) {
	return time.Time{}.Add(n.Duration).MarshalJSON()
}

func (n *NotificationAt) UnmarshalJSON(b []byte) error {
	t := time.Time{}
	if err := t.UnmarshalJSON(b); err != nil {
		return err
	}

	d, err := time.ParseDuration(t.Format("15h4m5s"))
	if err != nil {
		return err
	}
	*n = NotificationAt{Duration: d}

	return nil
}
func (n *NotificationAt) Scan(value any) error {
	td, err := time.ParseDuration(fmt.Sprintf("%dm", value.(int64)))
	if err != nil {
		return err
	}

	*n = NotificationAt{Duration: td}
	return nil
}

func (n NotificationAt) Value() (driver.Value, error) {
	return n.Duration.Minutes(), nil
}

func (n *NotificationAt) Parse(t time.Time) error {
	d, err := time.ParseDuration(t.Format("15h4m5s"))
	if err != nil {
		return err
	}
	*n = NotificationAt{Duration: d}
	return nil
}

type User struct {
	ID        int            `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email          string `gorm:"type:varchar(255)" json:"email"`
	SignInProvider string `gorm:"type:varchar(64)" json:"signInProvider"`
	Uid            string `gorm:"unique;type:varchar(64)" json:"-"`

	Role string `gorm:"type:varchar(64);default:'user'" json:"role"`

	Nickname *string `gorm:"type:varchar(64)" json:"nickname"`

	DeviceId *string `gorm:"type:varchar(36);uniqueIndex;" json:"-"`

	ProfilePicture         *string `gorm:"type:varchar(255)" json:"profilePicture"`
	ProfilePictureBlurHash *string `gorm:"type:varchar(30)" json:"profilePictureBlurHash"`

	FcmToken   *string `gorm:"type:varchar(255)" json:"-"`
	DeviceType *string `gorm:"type:varchar(32)" json:"-"`

	Notification   *bool           `json:"notification" gorm:"default:true"`
	NotificationAt *NotificationAt `json:"notificationAt" gorm:"type:int(4);default:0"`

	Locations []Location `json:"-"`
}

type UserUsecase interface {
	RegisterDevice(ctx context.Context, id int) (string, error)
	GetAuthUserByUid(ctx context.Context, uid string) (*User, error)
	VerifyIDToken(ctx context.Context, authorization string) (*User, error)
	GenerateIDToken(ctx context.Context, uid string) (string, error)
	Sync(ctx context.Context, user User) (*User, error)
	GetUserByDeviceId(ctx context.Context, deviceId string) (*User, error)
	ChangeProfile(ctx context.Context, file *multipart.FileHeader, hash string, userId int) error
	ChangeNickname(ctx context.Context, nickname string, userId int) error
	UpdateFcm(ctx context.Context, fcmToken *string, deviceType *string, userId int) error
	UpdateSettings(ctx context.Context, notification *bool, notificationsAt time.Time, userId int) error
	ListNotifications(ctx context.Context, notiAt time.Time) ([]*User, error)
}

type UserRepository interface {
	GetOrCreate(ctx context.Context, user User) (*User, error)
	Get(ctx context.Context, user User) (*User, error)
	UpdateByID(ctx context.Context, id int, user User) error
	ListNotifications(ctx context.Context, notiAt int) ([]*User, error)
}
