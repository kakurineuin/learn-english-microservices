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

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

const (
	WORD_MEANING_COLLECTION          = "wordmeanings"
	FAVORITE_WORD_MEANING_COLLECTION = "favoritewordmeanings"
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

func (repo *MongoDBRepository) CreateWordMeanings(
	ctx context.Context,
	wordMeanings []model.WordMeaning,
) (wordMeaningIds []string, err error) {
	now := time.Now()
	documents := []interface{}{}

	for i := range wordMeanings {
		wordMeanings[i].CreatedAt = now
		wordMeanings[i].UpdatedAt = now
		documents = append(documents, wordMeanings[i])
	}

	collection := repo.getCollection(WORD_MEANING_COLLECTION)
	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	for _, id := range result.InsertedIDs {
		wordMeaningIds = append(wordMeaningIds, id.(primitive.ObjectID).Hex())
	}

	return wordMeaningIds, nil
}

func (repo *MongoDBRepository) FindWordMeaningsByWordAndUserId(
	ctx context.Context,
	word, userId string,
) (wordMeanings []model.WordMeaning, err error) {
	matchStage := bson.D{{"$match", bson.D{{"queryByWords", word}}}}
	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "favoritewordmeanings"},
			{"localField", "_id"},
			{"foreignField", "wordMeaningId"},
			{"pipeline", bson.A{
				bson.D{{"$match", bson.D{{"userId", userId}}}},
			}},
			{"as", "favoriteWordMeanings"},
		},
	}}
	addFieldsStage := bson.D{{"$addFields", bson.D{
		{"favoriteWordMeaningId", bson.D{
			{"$cond", bson.A{
				bson.D{{"$gt", bson.A{
					bson.D{{"$size", "$favoriteWordMeanings"}},
					0,
				}}},
				bson.D{{"$arrayElemAt", bson.A{
					"$favoriteWordMeanings._id",
					0,
				}}},
				primitive.NilObjectID.Hex(), // 沒有找到 favoriteWordMeaning
			}},
		}},
	}}}
	projectStage := bson.D{{"$project", bson.D{{"favoriteWordMeanings", 0}}}}
	sortStage := bson.D{{"$sort", bson.D{{"orderByNo", 1}}}}

	wordMeaningsCollection := repo.getCollection(WORD_MEANING_COLLECTION)

	// pass the pipeline to the Aggregate() method
	cursor, err := wordMeaningsCollection.Aggregate(
		ctx,
		mongo.Pipeline{matchStage, lookupStage, addFieldsStage, projectStage, sortStage},
	)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &wordMeanings); err != nil {
		return nil, err
	}

	return wordMeanings, nil
}

func (repo *MongoDBRepository) CreateFavoriteWordMeaning(
	ctx context.Context,
	userId, wordMeaningId string,
) (favoriteWordMeaningId string, err error) {
	wordMeaningObjectId, err := primitive.ObjectIDFromHex(wordMeaningId)
	if err != nil {
		return "", err
	}

	colleciton := repo.getCollection(FAVORITE_WORD_MEANING_COLLECTION)
	result, err := colleciton.InsertOne(ctx, model.FavoriteWordMeaning{
		UserId:        userId,
		WordMeaningId: wordMeaningObjectId,
	})
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repo *MongoDBRepository) GetFavoriteWordMeaningById(
	ctx context.Context,
	favoriteWordMeaningId string,
) (favoriteWordMeaning *model.FavoriteWordMeaning, err error) {
	id, err := primitive.ObjectIDFromHex(favoriteWordMeaningId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{"_id", id},
	}
	var result model.FavoriteWordMeaning
	collection := repo.getCollection(FAVORITE_WORD_MEANING_COLLECTION)
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

func (repo *MongoDBRepository) FindFavoriteWordMeaningsByUserIdAndWord(
	ctx context.Context,
	userId, word string,
	skip, limit int64,
) (wordMeanings []model.WordMeaning, err error) {
	matchWord := bson.D{{}}

	if word != "" {
		matchWord = bson.D{{"wordMeaning.queryByWords", word}}
	}

	matchStage := bson.D{{"$match", bson.D{{"userId", userId}}}}
	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "wordmeanings"},
			{"localField", "wordMeaningId"},
			{"foreignField", "_id"},
			{"as", "wordMeaning"},
		},
	}}
	unwindStage := bson.D{{"$unwind", "$wordMeaning"}}
	matchWordStage := bson.D{{"$match", matchWord}}
	sortStage := bson.D{{"$sort", bson.D{{"wordMeaning.word", 1}, {"wordMeaning.orderByNo", 1}}}}
	skipStage := bson.D{{"$skip", skip}}
	limitStage := bson.D{{"$limit", limit}}

	collection := repo.getCollection(FAVORITE_WORD_MEANING_COLLECTION)

	// pass the pipeline to the Aggregate() method
	cursor, err := collection.Aggregate(
		ctx,
		mongo.Pipeline{
			matchStage,
			lookupStage,
			unwindStage,
			matchWordStage,
			sortStage,
			skipStage,
			limitStage,
		},
	)
	if err != nil {
		return nil, err
	}

	type result struct {
		Id            primitive.ObjectID `bson:"_id,omitempty"`
		UserId        string             `bson:"userId"`
		WordMeaningId primitive.ObjectID `bson:"wordMeaningId"`
		WordMeaning   model.WordMeaning  `bson:"wordMeaning"`
		CreatedAt     time.Time          `bson:"createdAt"`
		UpdatedAt     time.Time          `bson:"updatedAt"`
	}

	var results []result
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].WordMeaning.FavoriteWordMeaningId = results[i].Id
		wordMeanings = append(wordMeanings, results[i].WordMeaning)
	}

	return wordMeanings, nil
}

func (repo *MongoDBRepository) CountFavoriteWordMeaningsByUserIdAndWord(
	ctx context.Context,
	userId, word string,
) (count int64, err error) {
	matchWord := bson.D{{}}

	if word != "" {
		matchWord = bson.D{{"wordMeaning.queryByWords", word}}
	}

	matchStage := bson.D{{"$match", bson.D{{"userId", userId}}}}
	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "wordmeanings"},
			{"localField", "wordMeaningId"},
			{"foreignField", "_id"},
			{"as", "wordMeaning"},
		},
	}}
	unwindStage := bson.D{{"$unwind", "$wordMeaning"}}
	matchWordStage := bson.D{{"$match", matchWord}}
	countStage := bson.D{{"$count", "total"}}

	collection := repo.getCollection(FAVORITE_WORD_MEANING_COLLECTION)

	// pass the pipeline to the Aggregate() method
	cursor, err := collection.Aggregate(
		ctx,
		mongo.Pipeline{
			matchStage,
			lookupStage,
			unwindStage,
			matchWordStage,
			countStage,
		},
	)
	if err != nil {
		return 0, err
	}

	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	// const total = totalResult.length > 0 ? totalResult[0].total : 0;

	if len(results) > 0 {
		count = int64(results[0][0].Value.(int32))
	} else {
		count = 0
	}

	return count, nil
}

func (repo *MongoDBRepository) DeleteFavoriteWordMeaningById(
	ctx context.Context,
	favoriteWordMeaningId string,
) (deletedCount int64, err error) {
	id, err := primitive.ObjectIDFromHex(favoriteWordMeaningId)
	if err != nil {
		return 0, err
	}

	filter := bson.D{
		{"_id", id},
	}
	collection := repo.getCollection(FAVORITE_WORD_MEANING_COLLECTION)
	result, err := collection.DeleteOne(ctx, filter)
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
