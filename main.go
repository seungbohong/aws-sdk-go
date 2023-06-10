package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"s3test/constants"
	"s3test/handlers"
	"s3test/routes"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo 인스턴스 생성
	e := echo.New()

	// Middleware 등록
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// MongoDB 연결
	mongoClient, err := handlers.ConnectToMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// MongoDB 핸들러 생성
	db := handlers.NewMongoDB(mongoClient, constants.DBName, constants.CollectionName)

	// 라우트 초기화
	routes.InitRoutes(e, db)

	// 서버 시작
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 종료 시그널을 받을 채널 생성
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 서버 종료 대기 시간 설정 (2시간)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	// 서버 종료
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error: ", err)
	}
}
