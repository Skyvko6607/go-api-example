package main

import (
	"TestAPI/config"
	"TestAPI/database"
	"TestAPI/handlers"
	"TestAPI/repositories"
	"TestAPI/services"
	"context"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/pelletier/go-toml/v2"
)

func main() {
	r := gin.Default()

	configFile, configErr := os.ReadFile("appsettings.toml")
	if configErr != nil {
		panic(configErr)
	}

	var cfg config.AppSettings
	tomlErr := toml.Unmarshal([]byte(configFile), &cfg)
	if tomlErr != nil {
		panic(tomlErr)
	}

	mongoContext := database.NewMongoContext(cfg.MongoDatabase.Uri)

	defer func() {
		if err := mongoContext.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// User Service
	repo := &repositories.UserRepository{MongoContext: mongoContext}
	service := &services.UserService{Repo: repo}
	handler := &handlers.UserHandler{Service: service}

	// Index setup
	err := repo.EnsureUserIndexes(context.Background())
	if err != nil {
		panic(err)
	}

	r.GET("/users/:userNameOrEmail", handler.GetUser)
	r.POST("/users/create/:userName/:email", handler.CreateUser)
	r.Run(":8080")
}
