package main

import (
	"context"
	"os"

	"github.com/Skyvko6607/go-api-learning/config"
	"github.com/Skyvko6607/go-api-learning/database"
	"github.com/Skyvko6607/go-api-learning/handlers"
	"github.com/Skyvko6607/go-api-learning/repositories"
	"github.com/Skyvko6607/go-api-learning/services"

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
	err := toml.Unmarshal([]byte(configFile), &cfg)
	if err != nil {
		panic(err)
	}

	mongoContext := database.NewMongoContext(cfg.MongoDatabase.Uri)

	defer func() {
		if err := mongoContext.Client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// User Setup
	userRepo := &repositories.UserRepository{MongoContext: mongoContext}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}
	err = userRepo.EnsureIndexes(context.Background())
	if err != nil {
		panic(err)
	}
	userHandler.SetupEndpoints(r)

	// Product Setup
	productRepo := &repositories.ProductRepository{MongoContext: mongoContext}
	productService := &services.ProductService{Repo: productRepo}
	productHandler := &handlers.ProductHandler{Service: productService}
	productHandler.SetupEndpoints(r)

	// Order Setup
	orderRepo := &repositories.OrderRepository{MongoContext: mongoContext}
	orderService := &services.OrderService{Repo: orderRepo}
	orderHandler := &handlers.OrderHandler{Service: orderService}
	err = orderRepo.EnsureIndexes(context.Background())
	if err != nil {
		panic(err)
	}
	orderHandler.SetupEndpoints(r)

	r.Run(":8080")
}
