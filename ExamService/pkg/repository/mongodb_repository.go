package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

const (
	EXAM_COLLECTION        = "exams"
	QUESTION_COLLECTION    = "questions"
	ANSWERWRONG_COLLECTION = "answerwrongs"
	EXAM_RECORD_COLLECTION = "examrecords"
)

type MongoDBRepository struct {
	client   *mongo.Client
	database string
}

func NewMongoDBRepository(database string) *MongoDBRepository {
	return &MongoDBRepository{
		database: database,
	}
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

	collection := repo.getCollection(EXAM_COLLECTION)
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
	collection := repo.getCollection(EXAM_COLLECTION)
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

func (repo *MongoDBRepository) GetExamById(
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
	collection := repo.getCollection(EXAM_COLLECTION)
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

func (repo *MongoDBRepository) FindExamsByUserIdOrderByUpdateAtDesc(
	ctx context.Context,
	userId string,
	skip, limit int64,
) (exams []model.Exam, err error) {
	collection := repo.getCollection(EXAM_COLLECTION)
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

func (repo *MongoDBRepository) DeleteExamById(
	ctx context.Context,
	examId string,
) (deletedCount int64, err error) {
	id, err := primitive.ObjectIDFromHex(examId)
	if err != nil {
		return 0, err
	}

	filter := bson.D{
		{"_id", id},
	}
	collection := repo.getCollection(EXAM_COLLECTION)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (repo *MongoDBRepository) CountExamsByUserId(
	ctx context.Context,
	userId string,
) (count int64, err error) {
	collection := repo.getCollection(EXAM_COLLECTION)
	filter := bson.D{{"userId", userId}}
	count, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *MongoDBRepository) CreateQuestion(
	ctx context.Context,
	question model.Question,
) (questionId string, err error) {
	now := time.Now()
	question.CreatedAt = now
	question.UpdatedAt = now

	collection := repo.getCollection(QUESTION_COLLECTION)
	result, err := collection.InsertOne(ctx, question)
	if err != nil {
		return "", err
	}

	questionId = result.InsertedID.(primitive.ObjectID).Hex()
	return questionId, nil
}

func (repo *MongoDBRepository) UpdateQuestion(
	ctx context.Context,
	question model.Question,
) error {
	update := bson.D{{"$set", bson.D{
		{"examId", question.ExamId},
		{"ask", question.Ask},
		{"answers", question.Answers},
		{"userId", question.UserId},
		{"updatedAt", time.Now()},
	}}}
	collection := repo.getCollection(QUESTION_COLLECTION)
	result, err := collection.UpdateByID(ctx, question.Id, update)
	if err != nil {
		return err
	}

	// 查無符合條件的資料可供修改
	if result.MatchedCount == 0 {
		err = fmt.Errorf("Question not found by questionId: %s", question.Id.Hex())
		return err
	}

	return nil
}

func (repo *MongoDBRepository) GetQuestionById(
	ctx context.Context,
	questionId string,
) (question *model.Question, err error) {
	id, err := primitive.ObjectIDFromHex(questionId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{"_id", id},
	}
	var result model.Question
	collection := repo.getCollection(QUESTION_COLLECTION)
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

func (repo *MongoDBRepository) FindQuestionsByExamIdAndUserIdOrderByUpdateAtDesc(
	ctx context.Context,
	examId, userId string,
	skip, limit int64,
) (questions []model.Question, err error) {
	collection := repo.getCollection(QUESTION_COLLECTION)
	filter := bson.D{
		{"examId", examId},
		{"userId", userId},
	}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func (repo *MongoDBRepository) DeleteQuestionById(
	ctx context.Context,
	questionId string,
) (deletedCount int64, err error) {
	id, err := primitive.ObjectIDFromHex(questionId)
	if err != nil {
		return 0, err
	}

	filter := bson.D{
		{"_id", id},
	}
	collection := repo.getCollection(QUESTION_COLLECTION)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (repo *MongoDBRepository) DeleteQuestionsByExamId(
	ctx context.Context,
	examId string,
) (deletedCount int64, err error) {
	filter := bson.D{
		{"examId", examId},
	}
	collection := repo.getCollection(QUESTION_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (repo *MongoDBRepository) CountQuestionsByExamIdAndUserId(
	ctx context.Context,
	examId, userId string,
) (count int64, err error) {
	collection := repo.getCollection(QUESTION_COLLECTION)
	filter := bson.D{
		{"examId", examId},
		{"userId", userId},
	}
	count, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *MongoDBRepository) DeleteAnswerWrongByQuestionId(
	ctx context.Context,
	questionId string,
) (deletedCount int64, err error) {
	filter := bson.D{
		{"questionId", questionId},
	}
	collection := repo.getCollection(ANSWERWRONG_COLLECTION)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (repo *MongoDBRepository) DeleteAnswerWrongsByExamId(
	ctx context.Context,
	examId string,
) (deletedCount int64, err error) {
	filter := bson.D{
		{"examId", examId},
	}
	collection := repo.getCollection(ANSWERWRONG_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (repo *MongoDBRepository) DeleteExamRecordsByExamId(
	ctx context.Context,
	examId string,
) (deletedCount int64, err error) {
	filter := bson.D{
		{"examId", examId},
	}
	collection := repo.getCollection(EXAM_RECORD_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (repo *MongoDBRepository) WithTransaction(
	transactoinFunc transactionFunc,
) (interface{}, error) {
	// start-session
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := repo.client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("Start session failed! error: %w", err)
	}

	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(context.TODO())

	// Handle data within a transaction
	result, err := session.WithTransaction(
		context.TODO(),
		func(ctx mongo.SessionContext) (interface{}, error) {
			return transactoinFunc(ctx)
		},
		txnOptions,
	)
	return result, err
}

func (repo *MongoDBRepository) getCollection(collectionName string) *mongo.Collection {
	return repo.client.Database(repo.database).Collection(collectionName)
}
