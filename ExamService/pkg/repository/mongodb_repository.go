package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

const DATABASE = "learnEnglish"

type MongoDBRepository struct {
	client *mongo.Client
}

func NewMongoDBRepository() *MongoDBRepository {
	return &MongoDBRepository{}
}

func (repo *MongoDBRepository) ConnectDB(uri string) error {
	newClient, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("ConnectDB failed! error: %w", err)
	}

	repo.client = newClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ping the database
	err = repo.client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("ConnectDB ping database failed! error: %w", err)
	}

	fmt.Println("Connected to MongoDB")
	return nil
}

func (repo *MongoDBRepository) DisconnectDB() error {
	if err := repo.client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("DisconnectDB failed! error: %w", err)
	}

	return nil
}

func (repo *MongoDBRepository) CreateExam(
	ctx context.Context,
	exam model.Exam,
) (examId string, err error) {
	now := time.Now()
	exam.CreatedAt = now
	exam.UpdatedAt = now

	collection := repo.getCollection("exams")
	result, err := collection.InsertOne(ctx, exam)
	if err != nil {
		return "", err
	}

	examId = result.InsertedID.(primitive.ObjectID).Hex()
	return examId, nil
}

func (repo *MongoDBRepository) UpdateExam(
	ctx context.Context,
	exam model.Exam,
) error {
	update := bson.D{{"$set", bson.D{
		{"topic", exam.Topic},
		{"description", exam.Description},
		{"isPublic", exam.IsPublic},
		{"userId", exam.UserId},
		{"updatedAt", time.Now()},
	}}}
	collection := repo.getCollection("exams")
	result, err := collection.UpdateByID(ctx, exam.Id, update)
	if err != nil {
		return err
	}

	// 查無符合條件的資料可供修改
	if result.MatchedCount == 0 {
		err = fmt.Errorf("Exam not found by examId: %s", exam.Id.Hex())
		return err
	}

	return nil
}

func (repo *MongoDBRepository) GetExam(
	ctx context.Context,
	examId string,
) (exam *model.Exam, err error) {
	id, err := primitive.ObjectIDFromHex(examId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{"_id", id},
	}
	var result model.Exam
	collection := repo.getCollection("exams")
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		// 查無資料不視為錯誤
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (repo *MongoDBRepository) FindExamsOrderByUpdateAtDesc(
	ctx context.Context,
	userId string,
	skip, limit int64,
) (exams []model.Exam, err error) {
	collection := repo.getCollection("exams")
	filter := bson.D{{"userId", userId}}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &exams); err != nil {
		return nil, err
	}

	return exams, nil
}

func (repo *MongoDBRepository) DeleteExam(ctx context.Context, examId string) error {
	id, err := primitive.ObjectIDFromHex(examId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{"_id", id},
	}
	collection := repo.getCollection("exams")
	result, err := collection.DeleteOne(ctx, filter)

	// 查無符合條件的資料可供刪除
	if result.DeletedCount == 0 {
		err = fmt.Errorf("Exam not found by examId: %s", examId)
		return err
	}

	return nil
}

func (repo *MongoDBRepository) getCollection(collectionName string) *mongo.Collection {
	return repo.client.Database(DATABASE).Collection(collectionName)
}
