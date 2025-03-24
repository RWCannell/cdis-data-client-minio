package g3cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// go:generate mockgen -destination=./gen3-client/mocks/mock_gen3interface.go -package=mocks github.com/uc-cdis/gen3-client/gen3-client/g3cmd Gen3Interface

// useSSL for MinIO client
const useSSL = true

func InitMinIOClient(minioEndpoint string, minioAccessKeyID string, minioSecretAccessKey string) (*minio.Client, error) {
	// Initialize minio client object
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return minioClient, errors.New("Error has occurred while initializing MinIO client, detailed error message: " + err.Error())
	}

	return minioClient, err
}

func GenerateMinIOPresignedURL(minioClient *minio.Client, filename string) {
	// Set request parameters for content-disposition
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+filename+".txt\"")

	// Generates a presigned url which expires in a day.
	minioPresignedURL, err := minioClient.PresignedGetObject(context.Background(), "mybucket", "myobject", time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated MinIO presigned URL", minioPresignedURL)
}
