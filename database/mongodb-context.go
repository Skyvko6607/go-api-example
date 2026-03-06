package database

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoContext struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoContext(mongoDbUri string) *MongoContext {
	var client, err = mongo.Connect(options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		panic(err)
	}
	database := client.Database("testing")
	return &MongoContext{Client: client, Database: database}
}

func (m *MongoContext) Disconnect(c *gin.Context) {
	if err := m.Client.Disconnect(c); err != nil {
		panic(err)
	}
}
