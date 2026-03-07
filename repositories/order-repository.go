package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Skyvko6607/go-api-learning/database"
	"github.com/Skyvko6607/go-api-learning/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	MongoContext      *database.MongoContext
	ProductRepository *ProductRepository
}

func (r *OrderRepository) CreateOrder(c *gin.Context, createOrderDto models.CreateOrderDTO, userId bson.ObjectID) (models.Order, error) {
	collection := r.GetCollection()
	productsCollection := r.ProductRepository.GetCollection()
	filter := bson.M{
		"_id": bson.M{
			"$all": createOrderDto.Products,
		},
	}
	results, err := productsCollection.Find(c, filter)
	if err != nil {
		return models.Order{}, err
	}
	var products []models.Product
	err = results.All(c, &products)
	if err != nil {
		return models.Order{}, err
	}

	totalPaid := float32(0)
	for _, product := range products {
		totalPaid += product.Price
	}
	order := models.Order{
		UserID:      userId,
		ProductIDs:  createOrderDto.Products,
		Products:    products,
		CreatedAt:   time.Now(),
		PaymentType: createOrderDto.PaymentType,
		BillingData: createOrderDto.BillingData,
		TotalPaid:   totalPaid,
	}
	result, err := collection.InsertOne(c, order)
	if err != nil {
		return models.Order{}, err
	}
	if !result.Acknowledged {
		return models.Order{}, errors.New("failed to create order")
	}
	order.ID = result.InsertedID.(bson.ObjectID)
	return order, nil
}

func (r *OrderRepository) GetOrdersByUserId(c *gin.Context, userId bson.ObjectID) ([]models.Order, error) {
	collection := r.GetCollection()
	var orders []models.Order
	results, err := collection.Find(c, models.Order{
		UserID: userId,
	})
	if err != nil {
		return []models.Order{}, err
	}
	allErr := results.All(c, &orders)
	if allErr != nil {
		return []models.Order{}, allErr
	}
	return orders, nil
}

func (r *OrderRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.GetCollection()
	idx := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}
	if _, err := collection.Indexes().CreateOne(ctx, idx); err != nil {
		return err
	}
	idx = mongo.IndexModel{
		Keys: bson.D{{Key: "product_ids", Value: 1}},
	}
	_, err := collection.Indexes().CreateOne(ctx, idx)
	return err
}

func (r *OrderRepository) GetCollection() mongo.Collection {
	return *r.MongoContext.Database.Collection("orders")
}
