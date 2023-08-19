package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/V1taly5/APIYDisk/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	collection *mongo.Collection
}

func InitDataLayer() (*mongo.Client, error) {
	const op = "infrastructure.reposotory.mongodb_repository.InitDataLayer"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// URI move to env
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:rootpass@localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

func NewUserRepository(client *mongo.Client) entity.UserRepository {
	// nameDB & collName move to env
	db := client.Database("tgBot").Collection("user-token")

	return &userRepository{collection: db}
}

func (mongo *userRepository) FindByChatID(ctx context.Context, chatID int) (entity.User, error) {
	return mongo.findOneByQuery(ctx, bson.M{"chatid": chatID})
}

func (mongodb *userRepository) findOneByQuery(ctx context.Context, query interface{}) (entity.User, error) {
	const op = "infrastructure.reposotory.mongodb_repository.findOneByQuery"
	var result entity.User

	err := mongodb.collection.FindOne(ctx, query).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, mongo.ErrNoDocuments
		}
		return result, err
	}

	return result, nil
}

func (repo *userRepository) Create(ctx context.Context, user *entity.User) error {
	const op = "infrastructure.reposotory.mongodb_repository.Create"

	_, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (mongodb *userRepository) Update(ctx context.Context, user *entity.User) error {
	const op = "infrastructure.reposotory.mongodb_repository.Update"

	filter := bson.D{{Key: "chatid", Value: user.ChatID}}
	update := bson.D{{Key: "$set", Value: user}}

	_, err := mongodb.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
