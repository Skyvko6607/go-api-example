package handlers

import (
	"net/http"

	"github.com/Skyvko6607/go-api-example/models"
	"github.com/Skyvko6607/go-api-example/services"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service *services.OrderService
}

func (h *OrderHandler) SetupEndpoints(r *gin.Engine) {
	r.GET("/orders", h.GetOrders)
	r.POST("/orders", h.CreateOrder)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.Service.GetOrdersByUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var createOrderDto models.CreateOrderDTO
	if err := c.Bind(&createOrderDto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Service.CreateOrder(c, createOrderDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
