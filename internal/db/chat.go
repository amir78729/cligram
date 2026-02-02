package db

import (
	"cligram/internal/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepo struct {
	collection *mongo.Collection
}

func NewChatRepo(client *mongo.Client) *ChatRepo {
	coll := client.Database("cligram-db").Collection(string(ChatsCollection))

	// ensure unique index on "id"
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic("failed to create chat id index: " + err.Error())
	}

	return &ChatRepo{collection: coll}
}

// Create implements repository.ChatRepository
func (r *ChatRepo) Create(c domain.Chat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("chat with id %s already exists", c.ID)
		}
		return err
	}
	return nil
}

func (r *ChatRepo) GetByID(id string) (domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c domain.Chat
	err := r.collection.FindOne(ctx, map[string]string{"id": id}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Chat{}, domain.ErrChatNotFound
		}
		return domain.Chat{}, err
	}
	return c, nil
}
