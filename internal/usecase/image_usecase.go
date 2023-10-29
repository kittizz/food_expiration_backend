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
func (u *ImageUsecase) UploadImage(ctx context.Context, file *multipart.FileHeader, hash, folder string, userId int) (img *domain.Image, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	filename := uuid.New().String() + path.Ext(file.Filename)
	path := "/images/" + folder + "/" + filename

	_bytes, err := io.ReadAll(src)
	if err != nil {
		return
	}

	if hash == "" {
		buf := bytes.NewBuffer(_bytes)

		// Compute the BlurHash
		img, _, err := image.Decode(buf) // Note that image.Decode needs your file to be an image type it recognizes.
		if err != nil {
			return nil, err
		}

		hash, err = blurhash.Encode(4, 3, img)
		if err != nil {
			return nil, err
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
		return
	}

	img = &domain.Image{
		Path:     path,
		BlurHash: hash,
		UserID:   userId,
	}
	err = u.imageRepo.Create(ctx, img)

	if err != nil {
		return
	}

	return img, nil
}

func (u *ImageUsecase) Delete(ctx context.Context, id int) error {
	img, err := u.imageRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	err = u.bucket.RemoveObject(ctx, viper.GetString("BUCKET_NAME"), img.Path, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		return err
	}
	err = u.imageRepo.Delete(ctx, img.ID)
	if err != nil {
		return err
	}
	return nil
}
func (u *ImageUsecase) DeleteWithPath(ctx context.Context, path string) error {
	err := u.bucket.RemoveObject(ctx, viper.GetString("BUCKET_NAME"), path, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *ImageUsecase) Get(ctx context.Context, id int) (*domain.Image, error) {
	return u.imageRepo.Get(ctx, id)
}
