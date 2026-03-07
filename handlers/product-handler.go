package handlers

import (
	"net/http"

	"github.com/Skyvko6607/go-api-learning/models"
	"github.com/Skyvko6607/go-api-learning/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductHandler struct {
	Service *services.ProductService
}

func (h *ProductHandler) SetupEndpoints(r *gin.Engine) {
	r.GET("/products", h.GetAllProducts)
	r.POST("/products", h.AddProduct)
	r.DELETE("/products", h.DeleteProduct)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.Service.GetAllProducts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var productDto models.ProductDTO
	c.Bind(&productDto)

	product, err := h.Service.AddProduct(c, productDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productIdStr := c.Param("productId")
	productId, objErr := bson.ObjectIDFromHex(productIdStr)
	if objErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": objErr.Error()})
		return
	}

	err := h.Service.DeleteProduct(c, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
