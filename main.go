package main

import (
	"context"
	"os"

	"github.com/Skyvko6607/go-api-example/auth"
	"github.com/Skyvko6607/go-api-example/config"
	"github.com/Skyvko6607/go-api-example/database"
	"github.com/Skyvko6607/go-api-example/handlers"
	"github.com/Skyvko6607/go-api-example/repositories"
	"github.com/Skyvko6607/go-api-example/services"

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
	err := toml.Unmarshal(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	mongoContext := database.NewMongoContext(&cfg)
	redisContext := database.NewRedisContext(&cfg)

	defer func() {
		if err := mongoContext.Client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
		if err := redisContext.Client.Close(); err != nil {
			panic(err)
		}
	}()

	auth := &auth.Auth{RedisContext: &redisContext, AppSettings: &cfg}

	// User Setup
	userRepo := &repositories.UserRepository{MongoContext: mongoContext}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService, Auth: auth}
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

	authHandler := &handlers.AuthHandler{Auth: auth, Repo: userRepo, AppSettings: &cfg}
	authHandler.SetupEndpoints(r)

	r.Run(":8080")
}
