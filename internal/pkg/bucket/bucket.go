package bucket

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type Bucket struct {
	*minio.Client
}

func NewBucket() (*Bucket, error) {
	endpoint := viper.GetString("BUCKET_ENDPOINT")
	accessKeyID := viper.GetString("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := viper.GetString("BUCKET_SECRET_ACCESS_KEY")
	useSSL := viper.GetBool("BUCKET_USE_SSL")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	return &Bucket{minioClient}, err
}
