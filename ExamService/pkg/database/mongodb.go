package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/config"
)

var client *mongo.Client

func ConnectDB() error {
	newClient, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(config.EnvMongoDBURI()).SetTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("ConnectDB failed! error: %w", err)
	}

	client = newClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("ConnectDB ping database failed! error: %w", err)
	}

	fmt.Println("Connected to MongoDB")
	return nil
}

func DisconnectDB() error {
	if err := client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("DisconnectDB failed! error: %w", err)
	}

	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	collection := client.Database("learnEnglish").Collection(collectionName)
	return collection
}

type transactionFunc func(ctx mongo.SessionContext) (interface{}, error)

func WithTransaction(transactoinFunc transactionFunc) (interface{}, error) {
	// start-session
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("Start session failed! error: %w", err)
	}

	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(context.TODO())

	// Handle data within a transaction
	result, err := session.WithTransaction(context.TODO(), transactoinFunc, txnOptions)
	return result, err
}
