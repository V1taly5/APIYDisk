package repository

import (
	"context"
	"fmt"
	"log"
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

func (mongo *userRepository) FindByChatID(chatID int) (entity.User, error) {
	return mongo.findOneByQuery(bson.M{"chatid": chatID})
}

func (mongodb *userRepository) findOneByQuery(query interface{}) (entity.User, error) {
	const op = "infrastructure.reposotory.mongodb_repository.findOneByQuery"
	var result entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := mongodb.collection.FindOne(ctx, query).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, mongo.ErrNoDocuments
		}
		return result, err
	}

	return result, err
}

func (repo *userRepository) Create(user *entity.User) error {
	const op = "infrastructure.reposotory.mongodb_repository.Create"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (mongodb *userRepository) Update(user *entity.User) error {
	const op = "infrastructure.reposotory.mongodb_repository.Update"
	ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cencel()

	filter := bson.D{{Key: "chatid", Value: user.ChatID}}
	update := bson.D{{Key: "$set", Value: user}}

	_, err := mongodb.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (mongodb *userRepository) Find(chatID int) (bool, error) {
	var result entity.User
	ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cencel()
	err := mongodb.collection.FindOne(ctx, bson.M{"chatid": chatID}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		log.Fatal(err)
	}
	return true, nil
}

func (mongodb *userRepository) GetUser(chatID int64) (entity.User, error) {
	var result entity.User
	find, err := mongodb.Find(int(chatID))
	if err != nil {
		return result, err
	}
	if find {
		ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cencel()
		err = mongodb.collection.FindOne(ctx, bson.M{"chatid": chatID}).Decode(&result)
		if err != nil {
			return result, err
		}
		return result, nil
	} else {
		return result, mongo.ErrNilDocument
	}
}
