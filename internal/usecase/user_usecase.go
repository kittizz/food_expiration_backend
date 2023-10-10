package usecase

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"path"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/auth"
	"github.com/kittizz/food_expiration_backend/internal/pkg/bucket"
)

type UserUsecase struct {
	auth      *auth.Firebase
	userRepo  domain.UserRepository
	bucket    *bucket.Bucket
	imageRepo domain.ImageRepository
}

func NewUserUsecase(auth *auth.Firebase, userRepo domain.UserRepository, bucket *bucket.Bucket, imageRepo domain.ImageRepository) domain.UserUsecase {
	return &UserUsecase{
		auth:      auth,
		userRepo:  userRepo,
		bucket:    bucket,
		imageRepo: imageRepo,
	}
}

// GenerateIDToken implements domain.UserUsecase.
func (u *UserUsecase) GenerateIDToken(ctx context.Context, uid string) (string, error) {
	return u.auth.CustomToken(ctx, uid)
}

func (u *UserUsecase) VerifyIDToken(ctx context.Context, authorization string) (*domain.User, error) {

	token, err := u.auth.ParseIDToken(authorization)
	if err != nil {
		return nil, err
	}
	authToken, err := u.auth.VerifyIDToken(ctx, token)
	if err != nil {

		if e := log.Debug(); e.Enabled() {
			e.Msg(err.Error())

		}
		return nil, err
	}

	var email string
	if v, ok := authToken.Firebase.Identities["email"]; ok {
		email = v.([]any)[0].(string)
	}

	return &domain.User{
		Email:          email,
		SignInProvider: authToken.Firebase.SignInProvider,
		Uid:            authToken.UID,
	}, nil
}
func (u *UserUsecase) Sync(ctx context.Context, user domain.User) (*domain.User, error) {
	_user, err := u.userRepo.FetchOrCreate(ctx, user)

	return _user, err
}
func (u *UserUsecase) GetAuthUserByUid(ctx context.Context, uid string) (*domain.User, error) {
	_uid, err := u.auth.ParseIDToken(uid)
	if err != nil {
		return nil, err
	}
	userRecord, err := u.auth.GetUser(ctx, _uid)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Uid:            userRecord.UID,
		Email:          userRecord.Email,
		SignInProvider: userRecord.ProviderID,
	}, nil
}

// RegisterDevice implements domain.UserUsecase.
func (u *UserUsecase) RegisterDevice(ctx context.Context, id int) (string, error) {

	deviceId := uuid.New().String()
	err := u.userRepo.UpdateByID(ctx, id, domain.User{DeviceId: &deviceId})

	return deviceId, err
}

func (u *UserUsecase) GetUserByDeviceId(ctx context.Context, deviceId string) (*domain.User, error) {
	return u.userRepo.Fetch(ctx, domain.User{
		DeviceId: &deviceId,
	})
}

func (u *UserUsecase) ChangeProfile(ctx context.Context, file *multipart.FileHeader, userId int) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	filename := uuid.New().String() + path.Ext(file.Filename)
	path := "/images/" + filename

	_bytes, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(_bytes)

	_, err = u.bucket.PutObject(
		ctx,
		viper.GetString("BUCKET_NAME"),
		path,
		buf,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("content-type")},
	)
	if err != nil {
		return err
	}

	err = u.userRepo.UpdateByID(ctx, userId, domain.User{ProfilePicture: &path})
	if err != nil {
		return err
	}
	return nil
}
