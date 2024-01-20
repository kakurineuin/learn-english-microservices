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
	EXAM_COLLECTION         = "exams"
	QUESTION_COLLECTION     = "questions"
	ANSWER_WRONG_COLLECTION = "answerwrongs"
	EXAM_RECORD_COLLECTION  = "examrecords"
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

func (repo *MongoDBRepository) ConnectDB(ctx context.Context, uri string) error {
	newClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("ConnectDB failed! error: %w", err)
	}

	repo.client = newClient
	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ping the database
	err = repo.client.Ping(pingCtx, nil)
	if err != nil {
		return fmt.Errorf("ConnectDB ping database failed! error: %w", err)
	}

	fmt.Println("Connected to MongoDB")
	return nil
}

func (repo *MongoDBRepository) DisconnectDB(ctx context.Context) error {
	if err := repo.client.Disconnect(ctx); err != nil {
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
	skip, limit int32,
) (exams []model.Exam, err error) {
	collection := repo.getCollection(EXAM_COLLECTION)
	filter := bson.D{{"userId", userId}}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort).SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &exams); err != nil {
		return nil, err
	}

	return exams, nil
}

func (repo *MongoDBRepository) FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(
	ctx context.Context,
	userId string,
	isPublic bool,
) (exams []model.Exam, err error) {
	collection := repo.getCollection(EXAM_COLLECTION)
	filter := bson.D{
		{"userId", userId},
		{"isPublic", isPublic},
	}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort)
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
) (deletedCount int32, err error) {
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

	return int32(result.DeletedCount), nil
}

func (repo *MongoDBRepository) CountExamsByUserId(
	ctx context.Context,
	userId string,
) (count int32, err error) {
	collection := repo.getCollection(EXAM_COLLECTION)
	filter := bson.D{{"userId", userId}}
	result, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result), nil
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

func (repo *MongoDBRepository) FindQuestionsByQuestionIds(
	ctx context.Context,
	questionIds []string,
) (questions []model.Question, err error) {
	ids := []primitive.ObjectID{}

	for _, questionId := range questionIds {
		id, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	collection := repo.getCollection(QUESTION_COLLECTION)
	filter := bson.D{
		{"_id", bson.D{{"$in", ids}}},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func (repo *MongoDBRepository) FindQuestionsByExamIdOrderByUpdateAtDesc(
	ctx context.Context,
	examId string,
	skip, limit int32,
) (questions []model.Question, err error) {
	collection := repo.getCollection(QUESTION_COLLECTION)
	filter := bson.D{
		{"examId", examId},
	}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort).SetSkip(int64(skip)).SetLimit(int64(limit))
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
) (deletedCount int32, err error) {
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

	return int32(result.DeletedCount), nil
}

func (repo *MongoDBRepository) DeleteQuestionsByExamId(
	ctx context.Context,
	examId string,
) (deletedCount int32, err error) {
	filter := bson.D{
		{"examId", examId},
	}
	collection := repo.getCollection(QUESTION_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result.DeletedCount), nil
}

func (repo *MongoDBRepository) CountQuestionsByExamId(
	ctx context.Context,
	examId string,
) (count int32, err error) {
	collection := repo.getCollection(QUESTION_COLLECTION)
	filter := bson.D{
		{"examId", examId},
	}
	result, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result), nil
}

func (repo *MongoDBRepository) DeleteAnswerWrongsByQuestionId(
	ctx context.Context,
	questionId string,
) (deletedCount int32, err error) {
	filter := bson.D{
		{"questionId", questionId},
	}
	collection := repo.getCollection(ANSWER_WRONG_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result.DeletedCount), nil
}

func (repo *MongoDBRepository) DeleteAnswerWrongsByExamId(
	ctx context.Context,
	examId string,
) (deletedCount int32, err error) {
	filter := bson.D{
		{"examId", examId},
	}
	collection := repo.getCollection(ANSWER_WRONG_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result.DeletedCount), nil
}

func (repo *MongoDBRepository) UpsertAnswerWrongByTimesPlusOne(
	ctx context.Context,
	examId, questionId, userId string,
) (modifiedCount, upsertedCount int32, err error) {
	filter := bson.D{
		{"examId", examId},
		{"questionId", questionId},
		{"userId", userId},
	}

	now := time.Now()

	// times 遞增 1
	update := bson.D{
		{"$inc", bson.D{
			{"times", 1},
		}},
		{"$set", bson.D{
			{"updatedAt", now},
		}},
	}

	// Enable update or insert
	opts := options.Update().SetUpsert(true)

	collection := repo.getCollection(ANSWER_WRONG_COLLECTION)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return 0, 0, err
	}

	// 若是新增，補上新增日期時間
	if result.UpsertedCount == 1 {
		updateCreatedAt := bson.D{
			{"$set", bson.D{
				{"createdAt", now},
			}},
		}
		_, err := collection.UpdateOne(ctx, filter, updateCreatedAt)
		if err != nil {
			return 0, 0, err
		}
	}

	return int32(result.ModifiedCount), int32(result.UpsertedCount), nil
}

func (repo *MongoDBRepository) FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc(
	ctx context.Context,
	examId, userId string,
	limit int32,
) (answerWrongs []model.AnswerWrong, err error) {
	collection := repo.getCollection(ANSWER_WRONG_COLLECTION)
	filter := bson.D{
		{"examId", examId},
		{"userId", userId},
	}
	sort := bson.D{{"times", -1}} // descending
	opts := options.Find().SetSort(sort).SetLimit(int64(limit))
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &answerWrongs); err != nil {
		return nil, err
	}

	return answerWrongs, nil
}

func (repo *MongoDBRepository) DeleteExamRecordsByExamId(
	ctx context.Context,
	examId string,
) (deletedCount int32, err error) {
	filter := bson.D{
		{"examId", examId},
	}
	collection := repo.getCollection(EXAM_RECORD_COLLECTION)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result.DeletedCount), nil
}

func (repo *MongoDBRepository) CreateExamRecord(
	ctx context.Context,
	examRecord model.ExamRecord,
) (examRecordId string, err error) {
	now := time.Now()
	examRecord.CreatedAt = now
	examRecord.UpdatedAt = now

	collection := repo.getCollection(EXAM_RECORD_COLLECTION)
	result, err := collection.InsertOne(ctx, examRecord)
	if err != nil {
		return "", err
	}

	examRecordId = result.InsertedID.(primitive.ObjectID).Hex()
	return examRecordId, nil
}

func (repo *MongoDBRepository) FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
	ctx context.Context,
	examId, userId string,
	skip, limit int32,
) (examRecords []model.ExamRecord, err error) {
	collection := repo.getCollection(EXAM_RECORD_COLLECTION)
	filter := bson.D{
		{"examId", examId},
		{"userId", userId},
	}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort).SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &examRecords); err != nil {
		return nil, err
	}

	return examRecords, nil
}

func (repo *MongoDBRepository) CountExamRecordsByExamIdAndUserId(
	ctx context.Context,
	examId, userId string,
) (count int32, err error) {
	collection := repo.getCollection(EXAM_RECORD_COLLECTION)
	filter := bson.D{
		{"examId", examId},
		{"userId", userId},
	}
	result, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result), nil
}

func (repo *MongoDBRepository) FindExamRecordsByExamIdAndUserIdAndCreatedAt(
	ctx context.Context,
	examId,
	userId string,
	createdAt time.Time,
) (examRecords []model.ExamRecord, err error) {
	collection := repo.getCollection(EXAM_RECORD_COLLECTION)
	filter := bson.D{
		{"examId", examId},
		{"userId", userId},
		{"createdAt", bson.D{{"$gte", primitive.NewDateTimeFromTime(createdAt)}}},
	}
	sort := bson.D{{"createdAt", 1}} // ascending
	opts := options.Find().SetSort(sort)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &examRecords); err != nil {
		return nil, err
	}

	return examRecords, nil
}

func (repo *MongoDBRepository) WithTransaction(
	ctx context.Context,
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
	defer session.EndSession(ctx)

	// Handle data within a transaction
	result, err := session.WithTransaction(
		ctx,
		func(sctx mongo.SessionContext) (interface{}, error) {
			return transactoinFunc(sctx)
		},
		txnOptions,
	)
	return result, err
}

func (repo *MongoDBRepository) getCollection(collectionName string) *mongo.Collection {
	return repo.client.Database(repo.database).Collection(collectionName)
}
