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

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
)

const (
	USER_COLLECTION         = "users"
	USER_HISTORY_COLLECTION = "userhistories"
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

func (repo *MongoDBRepository) CreateUser(
	ctx context.Context,
	user model.User,
) (userId string, err error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	collection := repo.getCollection(USER_COLLECTION)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	userId = result.InsertedID.(primitive.ObjectID).Hex()
	return userId, nil
}

func (repo *MongoDBRepository) GetUserById(
	ctx context.Context,
	userId string,
) (user *model.User, err error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{"_id", id},
	}
	var result model.User
	collection := repo.getCollection(USER_COLLECTION)
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

func (repo *MongoDBRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (user *model.User, err error) {
	filter := bson.D{
		{"username", username},
	}
	var result model.User
	collection := repo.getCollection(USER_COLLECTION)
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

func (repo *MongoDBRepository) GetAdminUser(
	ctx context.Context,
) (user *model.User, err error) {
	filter := bson.D{
		{"role", "admin"},
	}
	var result model.User
	collection := repo.getCollection(USER_COLLECTION)
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

func (repo *MongoDBRepository) FindUserHistoryResponsesOrderByUpdatedAt(
	ctx context.Context, pageIndex, pageSize int32,
) (userHistoryResponses []UserHistoryResponse, err error) {
	sortStage := bson.D{{"$sort", bson.D{{"updatedAt", -1}}}}
	skipStage := bson.D{{"$skip", pageIndex * pageSize}}
	limitStage := bson.D{{"$limit", pageSize}}

	// 在 addFields 階段判斷 userId 是否為空字串
	addUserObjectIdStage := bson.D{{
		"$addFields", bson.D{
			{"userObjectId", bson.D{
				{"$cond", bson.A{
					bson.D{
						{"$eq", bson.A{"$userId", ""}},
					},
					nil,
					bson.D{
						{"$toObjectId", "$userId"},
					},
				}},
			}},
		},
	}}

	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "users"},
			{"localField", "userObjectId"},
			{"foreignField", "_id"},
			{"as", "user"},
		},
	}}
	unwindStage := bson.D{
		{"$unwind", bson.D{
			{"path", "$user"},
			{"preserveNullAndEmptyArrays", true},
		}},
	}
	addFieldsStage := bson.D{{
		"$addFields", bson.D{
			{
				"_id", bson.D{{"$toString", "$_id"}},
			},
			{
				"username", "$user.username",
			},
			{
				"role", "$user.role",
			},
		},
	}}
	projectStage := bson.D{{"$project", bson.D{{"userObjectId", 0}, {"user", 0}}}}

	collection := repo.getCollection(USER_HISTORY_COLLECTION)

	// pass the pipeline to the Aggregate() method
	cursor, err := collection.Aggregate(
		ctx,
		mongo.Pipeline{
			sortStage,
			skipStage,
			limitStage,
			addUserObjectIdStage,
			lookupStage,
			unwindStage,
			addFieldsStage,
			projectStage,
		},
	)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &userHistoryResponses); err != nil {
		return nil, err
	}

	// Empty slice 轉成前端 JSON 是空陣列
	// Nil slice 轉成前端 JSON 是 null
	// 方便前端不用判斷 null，所以使用 empty slice
	if userHistoryResponses == nil {
		userHistoryResponses = []UserHistoryResponse{}
	}

	return userHistoryResponses, nil
}

func (repo *MongoDBRepository) CountUserHistories(ctx context.Context) (count int32, err error) {
	collection := repo.getCollection(USER_HISTORY_COLLECTION)
	filter := bson.D{{}}
	result, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(result), nil
}

func (repo *MongoDBRepository) CreateUserHistory(
	ctx context.Context, userHistory model.UserHistory,
) (userHistoryId string, err error) {
	now := time.Now()
	userHistory.CreatedAt = now
	userHistory.UpdatedAt = now

	collection := repo.getCollection(USER_HISTORY_COLLECTION)
	result, err := collection.InsertOne(ctx, userHistory)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
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
