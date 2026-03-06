package repositories

import (
	"TestAPI/database"
	"TestAPI/models"
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	MongoContext *database.MongoContext
}

func (r *UserRepository) FindByUserNameOrEmail(userName string, email string) (models.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"userName": userName},
			{"email": email},
		},
	}
	collation := &options.Collation{Locale: "en", Strength: 2}

	return r.FindByBsonAndCollation(filter, collation)
}

func (r *UserRepository) FindByBsonAndCollation(filter bson.M, collation *options.Collation) (models.User, error) {
	collection := r.MongoContext.Database.Collection("users")

	var opts options.FindOneOptionsBuilder

	if collation != nil {
		opts = *options.FindOne().SetCollation(collation)
	} else {
		opts = *options.FindOne()
	}

	var u models.User
	err := collection.
		FindOne(context.TODO(), filter, &opts).
		Decode(&u)
	return u, err
}

func (r *UserRepository) CreateUser(userName string, email string) (models.User, error) {
	collection := r.MongoContext.Database.Collection("users")
	_, err := r.FindByUserNameOrEmail(userName, email)
	if err == nil {
		return models.User{}, errors.New("user already exists")
	}

	user := models.User{
		ID:            bson.NewObjectID(),
		UserName:      userName,
		UserNameLower: strings.ToLower(userName),
		Email:         email,
		EmailLower:    strings.ToLower(email),
	}
	_, err2 := collection.InsertOne(context.TODO(), user, options.InsertOne())
	return user, err2
}

func (r *UserRepository) EnsureUserIndexes(ctx context.Context) error {
	collection := r.MongoContext.Database.Collection("users")
	idx := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2,
		}),
	}
	if _, err := collection.Indexes().CreateOne(ctx, idx); err != nil {
		return err
	}

	idx2 := mongo.IndexModel{
		Keys: bson.D{{Key: "userName", Value: 1}},
		Options: options.Index().SetUnique(true).SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2,
		}),
	}
	_, err := collection.Indexes().CreateOne(ctx, idx2)
	return err
}
