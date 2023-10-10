package usecase

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path"

	"github.com/bbrks/go-blurhash"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/bucket"
)

type ImageUsecase struct {
	imageRepo domain.ImageRepository

	bucket *bucket.Bucket
}

func NewImageUsecase(imageRepo domain.ImageRepository, bucket *bucket.Bucket) domain.ImageUsecase {
	return &ImageUsecase{
		imageRepo: imageRepo,
		bucket:    bucket,
	}
}

// UploadImage implements domain.ImageUsecase.
func (u *ImageUsecase) UploadImage(ctx context.Context, file *multipart.FileHeader, userId int) error {
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

	// Compute the BlurHash
	img, _, err := image.Decode(buf) // Note that image.Decode needs your file to be an image type it recognizes.
	if err != nil {
		return err
	}

	hash, err := blurhash.Encode(4, 3, img)
	if err != nil {
		return err
	}

	buf = bytes.NewBuffer(_bytes)

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

	err = u.imageRepo.Create(ctx, &domain.Image{
		Path:     path,
		BlurHash: hash,
		UserID:   userId,
	})

	if err != nil {
		return err
	}

	return nil
}
