package upload

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioUpload() string {
	endpoint := "127.0.0.1:9001"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	filePath := "/Users/wyh/Macwork/mixed/1.png"
	bucketName := "alton"
	objectName := "c.png"

	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		log.Fatalln(errBucketExists)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Successfully created %s\n", bucketName)
	}
	contentType := "image/png"
	_, err = minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	//_, err = minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24*7, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(presignedURL)
	fmt.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, objectName)
	return presignedURL
}
