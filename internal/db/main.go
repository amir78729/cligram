package db

import (
	"cligram/internal/theme/printer"
	"context"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionName string

const (
	UsersCollection    CollectionName = "users"
	MessagesCollection CollectionName = "messages"
	ChatsCollection    CollectionName = "chats"
)

var (
	clientInstance *mongo.Client
	mongoOnce      sync.Once
)

func Connect() *mongo.Client {
	mongoOnce.Do(func() {
		_ = godotenv.Load(".env")
		uri := os.Getenv("MONGO_URI")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			printer.Error.Println("DB Connection Error:", err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			printer.Error.Println("DB Ping Error:", err)
		}

		clientInstance = client
		printer.Success.Println("Connected to DB")
	})
	return clientInstance
}

func GetCollection(name CollectionName) *mongo.Collection {
	return clientInstance.Database("cligram-db").Collection(string(name))
}
