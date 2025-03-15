package main

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "203.194.113.6:9000"
	accessKey := "MidFtK0wfiZ6AUjDfZbz"
	secretKey := "KxkgFNq196ok2AKq9U5h2naOUq0Akpi8HyjA4RO3"
	bucketName := "smk-telkom"

	fmt.Println("Connecting to MinIO...")
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		fmt.Println("Error connecting to MinIO:", err)
		return
	}

	fmt.Println("Connected to MinIO!")

	for {
		fmt.Println("\n== MinIO Program ==")
		fmt.Println("1. Upload file")
		fmt.Println("2. Download file")
		fmt.Println("3. List files")
		fmt.Println("4. Delete file")
		fmt.Println("5. Make presign URL")
		fmt.Println("6. Make public URL")
		fmt.Println("0. Exit")

		var choice int
		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			uploadFile(minioClient, bucketName)
		case 2:
			downloadFile(minioClient, bucketName)
		case 3:
			listFiles(minioClient, bucketName)
		case 4:
			deleteFile(minioClient, bucketName)
		case 5:
			makePresignURL(minioClient, bucketName)
		case 6:
			makePublicURL(minioClient, bucketName, endpoint)
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice!")
		}
	}
}

func uploadFile(minioClient *minio.Client, bucketName string) {
	var localFile string
	var objectName string

	fmt.Print("Enter file path: ")
	fmt.Scan(&localFile)

	fmt.Print("Enter object name: ")
	fmt.Scan(&objectName)

	_, err := os.Stat(localFile)
	if err != nil {
		fmt.Println("File not found!")
		return
	}

	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, localFile, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("Upload error:", err)
		return
	}

	fmt.Println("File uploaded successfully!")
}

func downloadFile(minioClient *minio.Client, bucketName string) {
	var objectName string
	var localFile string

	fmt.Print("Enter object name: ")
	fmt.Scan(&objectName)

	fmt.Print("Enter save path: ")
	fmt.Scan(&localFile)

	err := minioClient.FGetObject(context.Background(), bucketName, objectName, localFile, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Download error:", err)
		return
	}

	fmt.Println("File downloaded successfully!")
}

func listFiles(minioClient *minio.Client, bucketName string) {
	fmt.Println("Files in bucket:")

	objects := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{Recursive: true})

	for object := range objects {
		if object.Err != nil {
			fmt.Println("Error listing:", object.Err)
			return
		}
		fmt.Println("-", object.Key)
	}
}

func deleteFile(minioClient *minio.Client, bucketName string) {
	var objectName string

	fmt.Print("Enter object name to delete: ")
	fmt.Scan(&objectName)

	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Println("Delete error:", err)
		return
	}

	fmt.Println("File deleted successfully!")
}

func makePresignURL(minioClient *minio.Client, bucketName string) {
	var objectName string

	fmt.Print("Enter object name: ")
	fmt.Scan(&objectName)

	url, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, 3600, nil)
	if err != nil {
		fmt.Println("Error creating URL:", err)
		return
	}

	fmt.Println("Presigned URL:", url.String())
}

func makePublicURL(minioClient *minio.Client, bucketName string, endpoint string) {
	var objectName string
	var newObjectName string

	fmt.Print("Enter file path: ")
	fmt.Scan(&objectName)

	newObjectName = "/public/" + objectName

	fmt.Println("Public URL:", "http://"+endpoint+"/"+bucketName+newObjectName)
}
