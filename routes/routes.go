package routes

import (
	mongodb "s3test/handlers/mongo"
	"s3test/handlers/s3"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *mongodb.MongoDB, s3client *s3.S3Client, sess *session.Session) {
	// 라우트 설정
	e.GET("/upload", func(c echo.Context) error {
		return s3client.UploadFile(c, db)
	})
	e.GET("/download", func(c echo.Context) error {
		return s3client.DownloadFile(c, db, sess)
	})
	e.DELETE("/delete", func(c echo.Context) error {
		return s3client.DeleteFile(c, db)
	})
	e.PUT("/change-storage-class", func(c echo.Context) error {
		return s3client.ChangeStorageClass(c, db)
	})
	e.GET("/presigned-url", func(c echo.Context) error {
		return s3client.GetPresignedUrl(c)
	})
}
