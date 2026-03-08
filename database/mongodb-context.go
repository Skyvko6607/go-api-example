package database

import (
	"github.com/Skyvko6607/go-api-learning/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoContext struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoContext(appSettings *config.AppSettings) *MongoContext {
	var client, err = mongo.Connect(options.Client().ApplyURI(appSettings.MongoDatabase.Uri))
	if err != nil {
		panic(err)
	}
	database := client.Database(appSettings.MongoDatabase.Database)
	return &MongoContext{Client: client, Database: database}
}

func (m *MongoContext) Disconnect(c *gin.Context) {
	if err := m.Client.Disconnect(c); err != nil {
		panic(err)
	}
}
