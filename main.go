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

	configFile, configErr := os.ReadFile(".env")
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
	userRepo := &repositories.UserRepository{MongoContext: mongoContext}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}

	// Index setup
	userIndexErr := userRepo.EnsureIndexes(context.Background())
	if userIndexErr != nil {
		panic(userIndexErr)
	}

	productRepo := &repositories.ProductRepository{MongoContext: mongoContext}
	productService := &services.ProductService{Repo: productRepo}
	productHandler := &handlers.ProductHandler{Service: productService}

	userHandler.SetupEndpoints(r)
	productHandler.SetupEndpoints(r)
	r.Run(":8080")
}
