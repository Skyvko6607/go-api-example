package repositories

import (
	"errors"

	"github.com/Skyvko6607/go-api-learning/database"
	"github.com/Skyvko6607/go-api-learning/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository struct {
	MongoContext *database.MongoContext
}

func (r *ProductRepository) GetAllProducts(c *gin.Context) ([]models.Product, error) {
	collection := r.GetCollection()
	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return []models.Product{}, err
	}

	var products []models.Product
	allErr := cursor.All(c, &products)
	if allErr != nil {
		return []models.Product{}, err
	}

	return products, allErr
}

func (r *ProductRepository) GetProductById(c *gin.Context, productId bson.ObjectID) (models.Product, error) {
	collection := r.GetCollection()
	var product models.Product
	err := collection.FindOne(c, models.Product{
		ID: productId,
	}).Decode(&product)
	if err != nil {
		return models.Product{}, err
	}
	return product, err
}

func (r *ProductRepository) DeleteProduct(c *gin.Context, productId bson.ObjectID) error {
	collection := r.GetCollection()
	result, err := collection.DeleteOne(c, models.Product{
		ID: productId,
	})
	if err != nil {
		return err
	}

	if !result.Acknowledged {
		return errors.New("failed to remove product to database")
	}

	return nil
}

func (r *ProductRepository) AddProduct(c *gin.Context, product models.Product) (models.Product, error) {
	collection := r.GetCollection()
	result, err := collection.InsertOne(c, product)
	if err != nil {
		return models.Product{}, err
	}

	if !result.Acknowledged {
		return models.Product{}, errors.New("failed to add product to database")
	}

	product.ID = result.InsertedID.(bson.ObjectID)
	return product, nil
}

func (r *ProductRepository) GetCollection() *mongo.Collection {
	return r.MongoContext.Database.Collection("products")
}
