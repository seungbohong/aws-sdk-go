package mongodb

import (
	"context"
	"log"
	"s3test/constants"
	"s3test/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	client         *mongo.Client
	databaseName   string
	collectionName string
}

func NewMongoDB(client *mongo.Client, databaseName, collectionName string) *MongoDB {
	return &MongoDB{
		client:         client,
		databaseName:   databaseName,
		collectionName: collectionName,
	}
}

// 파일 메타데이터 삽입
func (db *MongoDB) InsertFileMetadata(file models.FileMetadata) error {
	// MongoDB 컬렉션 선택
	collection := db.client.Database(constants.DBName).Collection(constants.CollectionName)

	// FileMetadata 삽입
	_, err := collection.InsertOne(context.Background(), file)
	if err != nil {
		log.Println("Failed to insert file metadata to MongoDB", err)
		return err
	}

	log.Println("File metadata inserted to MongoDB")
	return nil
}

// 파일 메타데이터 조회
func (db *MongoDB) GetFileMetatData(fileName string) (*models.FileMetadata, error) {
	// MongoDB 컬렉션 선택
	collection := db.client.Database(constants.DBName).Collection(constants.CollectionName)

	// 조회할 파일의 조건 설정
	filter := bson.M{"filename": fileName}

	// 파일의 메타데이터 조회
	var fileMetadata models.FileMetadata
	err := collection.FindOne(context.Background(), filter).Decode(&fileMetadata)
	if err != nil {
		log.Println("Failed to get file metadata:", err)
		return nil, err
	}

	log.Println("Get File metadata from MongoDB:", &fileMetadata)
	return &fileMetadata, nil
}

// 파일 메타데이터 삭제
func (db *MongoDB) DeleteFileMeta(fileName string) error {
	// MongoDB 컬렉션 선택
	collection := db.client.Database(constants.DBName).Collection(constants.CollectionName)

	// 삭제할 파일의 조건 설정
	filter := bson.M{"filename": fileName}

	// 파일 메타데이터 삭제
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Failed to delete file metadata:", err)
		return err
	}

	log.Println("Delete File metadata from MongoDB:")
	return nil
}

func (db *MongoDB) UpdateFileMetadata(file models.FileMetadata) error {
	// MongoDB 컬렉션 선택
	collection := db.client.Database(constants.DBName).Collection(constants.CollectionName)

	// 업데이트할 파일의 조건 설정
	filter := bson.M{"filename": file.FileName}

	// 파일 메타데이터 업데이트 빌드 설정
	update := bson.D{
		{"$set", bson.D{{"storageclass", "INTELLIGENT_TIERING"}}},
	}

	// 파일 메타데이터 업데이트
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Failed to update file metadata:", err)
		return err
	}

	log.Println("Update the file metadata's storage class")
	return nil
}
