package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"s3test/constants"
	mongodb "s3test/handlers/mongo"
	s "s3test/handlers/s3"
	"s3test/routes"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Echo 인스턴스 생성
	e := echo.New()

	// Middleware 등록
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// MongoDB 연결
	mongoClient, err := ConnectToMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// MongoDB 인스턴스 생성
	db := mongodb.NewMongoDB(mongoClient, constants.DBName, constants.CollectionName)

	// S3 클라이언트 생성
	s3client, sess, _ := s.NewS3Client()

	// 라우트 초기화
	routes.InitRoutes(e, db, s3client, sess)

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

func ConnectToMongoDB() (*mongo.Client, error) {
	// MongoDB 연결 설정
	clientOptions := options.Client().
		ApplyURI(constants.MongoDBURL).
		SetAuth(options.Credential{
			// Get MongoDB credentials from Environment Variable
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		})
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB", err)
		return nil, err
	}

	// 연결 확인
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("Failed to ping MongoDB", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}
