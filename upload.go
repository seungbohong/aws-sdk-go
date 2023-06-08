package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	region   = "ap-northeast-1"
	bucket   = "seungbo-test"
	key      = "test.PNG"
	filePath = "./files/test.PNG"
)

func UploadFile() {
	// 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 파일 오픈
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 파일 업로드
	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		panic(err)
	}

	// 업로드 완료
	fmt.Println("파일 업로드 완료")
}

func DownloadFile() error {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return err
	}

	// S3 서비스 클라이언트 생성
	s3.New(sess)

	// 파일을 저장할 로컬 경로
	localFilePath := "downloaded-file.PNG"

	// 파일 다운로드를 위한 다운로더 생성
	downloader := s3manager.NewDownloader(sess)

	// localFilePath에 해당하는 파일 생성
	file, err := os.Create(localFilePath)
	if err != nil {
		fmt.Println("파일 생성 실패:", err)
		return err
	}
	defer file.Close()

	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("파일 다운로드 실패:", err)
		return err
	}

	fmt.Println("파일 다운로드 완료:", numBytes, "바이트")
	return nil
}

func DeleteFile() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region), // 필요한 리전으로 변경
	}))

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 객체 삭제 요청 생성
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// 객체 삭제
	_, err := svc.DeleteObject(deleteInput)
	if err != nil {
		fmt.Println("객체 삭제 실패:", err)
		return err
	}

	fmt.Println("객체 삭제 완료")
	return nil
}

func ChangeStorageClass() {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 객체의 스토리지 클래스 변경
	_, err = svc.CopyObject(&s3.CopyObjectInput{
		Bucket:            aws.String(bucket),
		CopySource:        aws.String(bucket + "/" + key),
		Key:               aws.String(key),
		StorageClass:      aws.String("INTELLIGENT_TIERING"),
		MetadataDirective: aws.String("COPY"),
	})
	if err != nil {
		panic(err)
	}

	// 스토리지 클래스 변경 완료
	fmt.Println("객체의 스토리지 클래스 변경 완료")
}

func GetPresignedUrl() {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// PresignedURL 생성 요청 설정
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// Presinged URL 만료 시간 설정 (예: 1시간)
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	log.Println("The URL is", urlStr)
}

func main() {
	// ChangeStorageClass()
	// fmt.Println("Successfully changed storage class")

	// GetPresignedUrl()

	// DownloadFile()

	DeleteFile()
}
