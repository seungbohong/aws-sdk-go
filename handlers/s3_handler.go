package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"s3test/constants"
	"s3test/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {
	// 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
		return err
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 파일 오픈
	file, err := os.Open(constants.FilePath)
	if err != nil {
		log.Println("Failed to open a file", err)
		return err
	}
	defer file.Close()

	// 파일 업로드
	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.FileName),
	})
	if err != nil {
		log.Println("Failed to upload a file", err)
		return err
	}

	// 파일 메타데이터 저장
	fileMeta := models.FileMetadata{
		FileName:   constants.FileName,
		FilePath:   constants.FilePath,
		UploadTime: time.Now(),
	}
	err = db.InsertFileMetadata(&fileMeta)
	if err != nil {
		log.Println("Failed to insert file metadata to MongoDB", err)
		return err
	}

	// 업로드 완료
	log.Println("파일 업로드 완료")
	return c.String(http.StatusOK, "File uploaded successfully")
}

func DownloadFile(c echo.Context) error {
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
		return err
	}
	defer file.Close()

	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.FileName),
	})
	if err != nil {
		log.Println("Failed to download a file", err)
		return err
	}

	// 파일 다운로드 후 메타데이터 조회
	fileMeta, err := db.GetFileMetatData(constants.FileName)
	if err != nil {
		log.Println("Failed to get file metadata from MongoDB", err)
		return err
	}

	log.Println("파일 다운로드 완료:", fileMeta)
	return c.String(http.StatusOK, fmt.Sprintf("File downloaded successfully: %d bytes", numBytes))
}

func DeleteFile(c echo.Context) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(constants.Region), // 필요한 리전으로 변경
	}))

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 파일 삭제 요청 생성
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.FileName),
	}

	// 파일 삭제
	_, err := svc.DeleteObject(deleteInput)
	if err != nil {
		log.Println("Failed to delete an object", err)
		return err
	}

	// 파일 삭제 후 메타데이터 삭제
	err = db.DeleteFileMeta(constants.FileName)
	if err != nil {
		log.Println("Failed to delete file metadata", err)
		return err
	}

	log.Println("객체 삭제 완료")
	return c.String(http.StatusOK, "File deleted successfully")
}

func ChangeStorageClass(c echo.Context) error {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
		return err
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// 객체의 스토리지 클래스 변경
	_, err = svc.CopyObject(&s3.CopyObjectInput{
		Bucket:            aws.String(constants.Bucket),
		CopySource:        aws.String(constants.Bucket + "/" + constants.FileName),
		Key:               aws.String(constants.FileName),
		StorageClass:      aws.String("INTELLIGENT_TIERING"),
		MetadataDirective: aws.String("COPY"),
	})
	if err != nil {
		log.Println("Failed to change an object's storage class", err)
		return err
	}

	err = db.UpdateFileMetadata(models.FileMetadata{
		FileName:     constants.FileName,
		StorageClass: "INTELLIGENT_TIERING",
	})
	if err != nil {
		log.Println("Failed to update file metadata in MongoDB", err)
		return err
	}

	// 스토리지 클래스 변경 완료
	log.Println("객체의 스토리지 클래스 변경 완료")
	return c.String(http.StatusOK, "Object's storage class changed successfully")
}

func GetPresignedUrl(c echo.Context) error {
	// AWS 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Println("Failed to create a session", err)
		return err
	}

	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// PresignedURL 생성 요청 설정
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(constants.FileName),
	})

	// Presinged URL 만료 시간 설정 (예: 1시간)
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
		return err
	}

	log.Println("The URL is", urlStr)
	return c.String(http.StatusOK, urlStr)
}
