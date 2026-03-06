package services

import (
	"TestAPI/models"
	"TestAPI/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

func (r *ProductService) GetAllProducts() ([]models.ProductDTO, error) {
	products, err := r.Repo.GetAllProducts()
	if err != nil {
		return []models.ProductDTO{}, err
	}
	var productsDTO []models.ProductDTO
	for _, product := range products {
		productsDTO = append(productsDTO, product.AsDTO())
	}

	return productsDTO, nil
}

func (r *ProductService) AddProduct(productDto models.ProductDTO) (models.ProductDTO, error) {
	product, err := r.Repo.AddProduct(productDto.AsBase())
	if err != nil {
		return models.ProductDTO{}, err
	}

	return product.AsDTO(), nil
}

func (r *ProductService) DeleteProduct(productId bson.ObjectID) error {
	return r.Repo.DeleteProduct(productId)
}
