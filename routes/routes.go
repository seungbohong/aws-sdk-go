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

	/*
		// MongoDB 관련 endpoints
		// 파일 메타데이터 조회
		e.GET("/file/:id", func(c echo.Context) error {
			fileName := c.Param("id")
			file, err := db.GetFileMetatData(fileName)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, file)
		})

		// 파일 메타데이터 삭제
		e.DELETE("/file/:id", func(c echo.Context) error {
			fileName := c.Param("id")
			err := db.DeleteFileMeta(fileName)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, fileName)
		})

		// 파일 메타데이터 업데이트
		e.PUT("/file/:id", func(c echo.Context) error {
			fileName := c.Param("id")
			var updatedFile models.FileMetadata
			if err := c.Bind(&updatedFile); err != nil {
				return err
			}
			updatedFile.FileName = fileName
			err := db.UpdateFileMetadata(updatedFile)
			if err != nil {
				return err
			}
			return c.String(http.StatusOK, "File metadata updated successfully")
		})
	*/
}
