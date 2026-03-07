package services

import (
	"github.com/Skyvko6607/go-api-learning/models"
	"github.com/Skyvko6607/go-api-learning/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderService struct {
	Repo *repositories.OrderRepository
}

func (s *OrderService) CreateOrder(c *gin.Context, createOrderDto models.CreateOrderDTO) (models.OrderDTO, error) {
	var userId bson.ObjectID // Get User id from auth
	order, err := s.Repo.CreateOrder(c, createOrderDto, userId)
	if err != nil {
		return models.OrderDTO{}, err
	}

	return order.AsDTO(), nil
}

func (s *OrderService) GetOrdersByUserId(c *gin.Context) ([]models.OrderDTO, error) {
	var userId bson.ObjectID // Get User id from auth
	orders, err := s.Repo.GetOrdersByUserId(c, userId)
	if err != nil {
		return []models.OrderDTO{}, err
	}
	var orderDtos []models.OrderDTO
	for _, order := range orders {
		orderDtos = append(orderDtos, order.AsDTO())
	}
	return orderDtos, nil
}
