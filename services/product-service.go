package services

import (
	"github.com/Skyvko6607/go-api-example/models"
	"github.com/Skyvko6607/go-api-example/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

func (r *ProductService) GetAllProducts(c *gin.Context) ([]models.ProductDTO, error) {
	products, err := r.Repo.GetAllProducts(c)
	if err != nil {
		return []models.ProductDTO{}, err
	}
	var productsDTO []models.ProductDTO
	for _, product := range products {
		productsDTO = append(productsDTO, product.AsDTO())
	}

	return productsDTO, nil
}

func (r *ProductService) AddProduct(c *gin.Context, productDto models.ProductDTO) (models.ProductDTO, error) {
	product, err := r.Repo.AddProduct(c, productDto.AsBase())
	if err != nil {
		return models.ProductDTO{}, err
	}

	return product.AsDTO(), nil
}

func (r *ProductService) DeleteProduct(c *gin.Context, productId bson.ObjectID) error {
	return r.Repo.DeleteProduct(c, productId)
}
