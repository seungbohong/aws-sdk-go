package routes

import (
	"s3test/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/upload", handler.UploadFile)
	e.GET("/download", handler.DownloadFile)
	e.GET("/delete", handler.DeleteFile)
	e.GET("/changeStorageClass", handler.ChangeStorageClass)
	e.GET("/GetPresignedUrl", handler.GetPresignedUrl)
}
