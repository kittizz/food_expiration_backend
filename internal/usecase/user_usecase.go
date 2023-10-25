package usecase

import (
	"bytes"
	"context"
	"image"
	"io"
	"mime/multipart"
	"path"

	"github.com/bbrks/go-blurhash"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/bucket"
	"github.com/kittizz/food_expiration_backend/internal/pkg/firebase"
)

type UserUsecase struct {
	firebase  *firebase.Firebase
	userRepo  domain.UserRepository
	bucket    *bucket.Bucket
	imageRepo domain.ImageRepository
}

func NewUserUsecase(firebase *firebase.Firebase, userRepo domain.UserRepository, bucket *bucket.Bucket, imageRepo domain.ImageRepository) domain.UserUsecase {
	return &UserUsecase{
		firebase:  firebase,
		userRepo:  userRepo,
		bucket:    bucket,
		imageRepo: imageRepo,
	}
}

// GenerateIDToken implements domain.UserUsecase.
func (u *UserUsecase) GenerateIDToken(ctx context.Context, uid string) (string, error) {
	return u.firebase.CustomToken(ctx, uid)
}

func (u *UserUsecase) VerifyIDToken(ctx context.Context, authorization string) (*domain.User, error) {

	token, err := u.firebase.ParseIDToken(authorization)
	if err != nil {
		return nil, err
	}
	authToken, err := u.firebase.VerifyIDToken(ctx, token)
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
	_user, err := u.userRepo.GetOrCreate(ctx, user)

	return _user, err
}
func (u *UserUsecase) GetAuthUserByUid(ctx context.Context, uid string) (*domain.User, error) {
	_uid, err := u.firebase.ParseIDToken(uid)
	if err != nil {
		return nil, err
	}
	userRecord, err := u.firebase.GetUser(ctx, _uid)
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
	return u.userRepo.Get(ctx, domain.User{
		DeviceId: &deviceId,
	})
}

func (u *UserUsecase) ChangeProfile(ctx context.Context, file *multipart.FileHeader, hash string, userId int) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	filename := uuid.New().String() + path.Ext(file.Filename)
	path := "/images/profile-picture/" + filename

	_bytes, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	if hash == "" {
		buf := bytes.NewBuffer(_bytes)

		// Compute the BlurHash
		img, _, err := image.Decode(buf) // Note that image.Decode needs your file to be an image type it recognizes.
		if err != nil {
			return err
		}

		hash, err = blurhash.Encode(4, 3, img)
		if err != nil {
			return err
		}

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

	err = u.userRepo.UpdateByID(ctx, userId, domain.User{ProfilePicture: &path, ProfilePictureBlurHash: &hash})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) ChangeNickname(ctx context.Context, nickname string, userId int) error {
	err := u.userRepo.UpdateByID(ctx, userId, domain.User{Nickname: &nickname})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) UpdateFcm(ctx context.Context, fcmToken *string, deviceType *string, userId int) error {
	err := u.userRepo.UpdateByID(ctx, userId, domain.User{FcmToken: fcmToken, DeviceType: deviceType})
	if err != nil {
		return err
	}
	return nil
}
