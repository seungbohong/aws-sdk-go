package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"s3test/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadFile() error {
	// 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 파일 오픈
	file, err := os.Open(constants.FilePath)
	if err != nil {
		log.Println("Failed to open a file", err)
	}
	defer file.Close()

	// 파일 업로드
	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.Key),
	})
	if err != nil {
		log.Println("Failed to upload a file", err)
	}

	// 업로드 완료
	fmt.Println("파일 업로드 완료")
	return nil
}

func DownloadFile() error {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
	}

	// S3 서비스 클라이언트 생성
	s3.New(sess)

	// 파일을 저장할 로컬 경로
	localFilePath := "downloadfiles/downloaded-file.PNG"

	// 파일 다운로드를 위한 다운로더 생성
	downloader := s3manager.NewDownloader(sess)

	// localFilePath에 해당하는 파일 생성
	file, err := os.Create(localFilePath)
	if err != nil {
		log.Println("Failed to create a file", err)
	}
	defer file.Close()

	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.Key),
	})
	if err != nil {
		log.Println("Failed to download a file", err)
		return err
	}

	fmt.Println("파일 다운로드 완료:", numBytes, "바이트")
	return nil
}

func DeleteFile() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(constants.Region), // 필요한 리전으로 변경
	}))

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 객체 삭제 요청 생성
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.Key),
	}

	// 객체 삭제
	_, err := svc.DeleteObject(deleteInput)
	if err != nil {
		log.Println("Failed to delete an object", err)
	}

	fmt.Println("객체 삭제 완료")
	return nil
}

func ChangeStorageClass() error {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 객체의 스토리지 클래스 변경
	_, err = svc.CopyObject(&s3.CopyObjectInput{
		Bucket:            aws.String(constants.Bucket),
		CopySource:        aws.String(constants.Bucket + "/" + constants.Key),
		Key:               aws.String(constants.Key),
		StorageClass:      aws.String("INTELLIGENT_TIERING"),
		MetadataDirective: aws.String("COPY"),
	})
	if err != nil {
		log.Println("Failed to change an object's storage class", err)
	}

	// 스토리지 클래스 변경 완료
	fmt.Println("객체의 스토리지 클래스 변경 완료")
	return nil
}

func GetPresignedUrl() error {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// PresignedURL 생성 요청 설정
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.Key),
	})

	// Presinged URL 만료 시간 설정 (예: 1시간)
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	log.Println("The URL is", urlStr)
	return nil
}

func main() {
	// s3 버킷에 파일 업로드
	// UploadFile()

	// 객체의 스토리지 클래스 변경
	// ChangeStorageClass()

	// s3 버킷에 있는객체의 Presigned URL
	// GetPresignedUrl()

	// s3 버킷에서 파일 다운로드
	// DownloadFile()

	// s3 버킷에서 객체 삭제
	DeleteFile()
}
