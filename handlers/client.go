package handlers

import (

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"kluisz-object-storage/config"
)


func getMinioClient() (*minio.Client, error) {
	return minio.New(config.Cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Cfg.S3.AccessKey, config.Cfg.S3.SecretKey, ""),
		Secure: config.Cfg.S3.UseSSL,
		Region: config.Cfg.S3.Region,
	})
}