package routes

import (
	"net/http"
	mongodb "s3test/handlers/mongo"
	"s3test/handlers/s3"
	"s3test/models"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *mongodb.MongoDB) {

	// 라우트 설정
	e.POST("/upload", s3.UploadFile)
	e.GET("/download", s3.DownloadFile)
	e.DELETE("/delete", s3.DeleteFile)
	e.PUT("/change-storage-class", s3.ChangeStorageClass)
	e.GET("/presigned-url", s3.GetPresignedUrl)

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
}
