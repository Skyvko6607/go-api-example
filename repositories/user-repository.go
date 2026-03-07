package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/Skyvko6607/go-api-learning/database"
	"github.com/Skyvko6607/go-api-learning/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	MongoContext *database.MongoContext
}

func (r *UserRepository) FindByUserNameOrEmail(c *gin.Context, userName string, email string) (models.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"user_name_lower": strings.ToLower(userName)},
			{"email": email},
		},
	}
	collation := &options.Collation{Locale: "en", Strength: 2}

	return r.FindByBsonAndCollation(c, filter, collation)
}

func (r *UserRepository) FindByBsonAndCollation(c *gin.Context, filter bson.M, collation *options.Collation) (models.User, error) {
	collection := r.MongoContext.Database.Collection("users")

	var opts options.FindOneOptionsBuilder

	if collation != nil {
		opts = *options.FindOne().SetCollation(collation)
	} else {
		opts = *options.FindOne()
	}

	var u models.User
	err := collection.
		FindOne(c, filter, &opts).
		Decode(&u)
	return u, err
}

func (r *UserRepository) CreateUser(c *gin.Context, userName string, email string, passwordHash string) (models.User, error) {
	collection := r.MongoContext.Database.Collection("users")
	_, err := r.FindByUserNameOrEmail(c, userName, email)
	if err == nil {
		return models.User{}, errors.New("user already exists")
	}

	user := models.User{
		ID:            bson.NewObjectID(),
		UserName:      userName,
		UserNameLower: strings.ToLower(userName),
		Email:         strings.ToLower(email),
		PasswordHash:  passwordHash,
	}
	_, err2 := collection.InsertOne(c, user, options.InsertOne())
	return user, err2
}

func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
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

	idx = mongo.IndexModel{
		Keys: bson.D{{Key: "user_name_lower", Value: 1}},
		Options: options.Index().SetUnique(true).SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2,
		}),
	}
	_, err := collection.Indexes().CreateOne(ctx, idx)
	return err
}
