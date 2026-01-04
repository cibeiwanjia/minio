package upload

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func RustFSUpload(endPoint, access, secret, bucket, object, fileRoad string) {
	endpoint := endPoint // RustFS 的 API 地址
	accessKey := access  // 你的 Access Key
	secretKey := secret  // 你的 Secret Key
	bucketName := bucket // 目标存储桶名称
	objectKey := object  // 文件在桶中的路径/名称
	filePath := fileRoad // 本地图片文件路径

	// 2. 创建自定义配置解析器，用于连接 RustFS
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: endpoint, // 指定 RustFS 的自定义端点
				}, nil
			}),
		),
		// 强制使用路径样式，这是连接非 AWS S3 服务的常见设置
		config.WithClientLogMode(aws.LogSigning),
	)
	if err != nil {
		log.Fatal("配置加载失败: ", err)
	}

	// 3. 创建 S3 客户端
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // 非常重要！设置为路径样式访问
	})

	// 4. 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("无法打开文件: ", err)
	}
	defer file.Close()

	// 5. 执行上传操作
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
		// 可以根据需要设置 Content-Type，例如 "image/jpeg"
		ContentType: aws.String("image/jpeg"),
		ACL:         "public-read", // 设置此文件的访问控制列表为公开读
	})
	if err != nil {
		log.Fatal("上传失败: ", err)
	}

	fmt.Printf("图片上传成功！地址: %s/%s\n", endpoint, bucketName+"/"+objectKey)
}
