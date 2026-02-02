package db

import (
	"cligram/internal/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepo struct {
	collection *mongo.Collection
}

func NewMessageRepo(client *mongo.Client) *MessageRepo {
	coll := client.Database("cligram-db").Collection(string(MessagesCollection))

	// ensure unique index on "id"
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic("failed to create message id index: " + err.Error())
	}

	return &MessageRepo{collection: coll}
}

// Create implements repository.MessageRepository
func (r *MessageRepo) Create(m domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, m)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("message with id %s already exists", m.ID)
		}
		return err
	}
	return nil
}

func (r *MessageRepo) ListByChat(chatID string) ([]domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, map[string]string{"chatid": chatID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []domain.Message
	for cursor.Next(ctx) {
		var m domain.Message
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
