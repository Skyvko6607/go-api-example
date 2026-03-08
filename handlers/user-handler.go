package handlers

import (
	"net/http"

	"github.com/Skyvko6607/go-api-learning/auth"
	"github.com/Skyvko6607/go-api-learning/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
	Auth    *auth.Auth
}

func (h *UserHandler) SetupEndpoints(e *gin.Engine) {
	e.GET("/users/:userNameOrEmail", h.GetUser)
	e.POST("/users/:userName/:email/:password", h.CreateUser)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userNameOrEmail := c.Param("userNameOrEmail")

	user, err := h.Service.GetUser(c, userNameOrEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	userName := c.Param("userName")
	email := c.Param("email")
	password := c.Param("password")

	user, err := h.Service.CreateUser(c, userName, email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
