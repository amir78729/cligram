package db

import (
	"cligram/internal/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(client *mongo.Client) *UserRepo {
	coll := client.Database("cligram-db").Collection(string(UsersCollection))

	// ensure unique index on "id"
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"id": 1}, // ascending
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic("failed to create user id index: " + err.Error())
	}

	return &UserRepo{collection: coll}
}

// Create implements repository.UserRepository
func (r *UserRepo) Create(u domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("user with id %s already exists", u.ID)
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByID(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u domain.User
	err := r.collection.FindOne(ctx, map[string]string{"id": id}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}
