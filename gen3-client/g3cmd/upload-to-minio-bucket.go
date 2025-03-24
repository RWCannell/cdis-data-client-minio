package g3cmd

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func uploadToMinIOBucket() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	endpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := true

	s3Client, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
			Secure: useSSL,
		})
	if err != nil {
		log.Fatalln(err)
	}

	uploadFileName := "If_by_Rudyard_Kipling.pdf"
	uploadFilePath := "../data_uploads/If_by_Rudyard_Kipling.pdf"

	if _, err := s3Client.FPutObject(context.Background(), bucketName, uploadFileName, uploadFilePath, minio.PutObjectOptions{
		ContentType: "application/*",
	}); err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully uploaded")
}
