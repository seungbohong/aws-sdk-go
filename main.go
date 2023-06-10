package main

import (
	"s3test/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Echo 인스턴스 생성
	e := echo.New()

	// 라우터 설정
	routes.RegisterRoutes(e)

	// 서버 실행
	e.Start(":8080")
}
