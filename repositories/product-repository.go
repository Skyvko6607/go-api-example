package repositories

import (
	"TestAPI/database"
	"TestAPI/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository struct {
	MongoContext *database.MongoContext
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	collection := r.GetCollection()
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return []models.Product{}, err
	}

	var products []models.Product
	allErr := cursor.All(context.TODO(), &products)
	if allErr != nil {
		return []models.Product{}, err
	}

	return products, allErr
}

func (r *ProductRepository) GetProductById(productId bson.ObjectID) (models.Product, error) {
	collection := r.GetCollection()
	var product models.Product
	err := collection.FindOne(context.TODO(), models.Product{
		ID: productId,
	}).Decode(&product)
	if err != nil {
		return models.Product{}, err
	}
	return product, err
}

func (r *ProductRepository) DeleteProduct(productId bson.ObjectID) error {
	collection := r.GetCollection()
	result, err := collection.DeleteOne(context.TODO(), models.Product{
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

func (r *ProductRepository) AddProduct(product models.Product) (models.Product, error) {
	collection := r.GetCollection()
	result, err := collection.InsertOne(context.TODO(), product)
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
