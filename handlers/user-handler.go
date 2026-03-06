package handlers

import (
	"TestAPI/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userNameOrEmail := c.Param("userNameOrEmail")

	user, err := h.Service.GetUser(userNameOrEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	userName := c.Param("userName")
	email := c.Param("email")

	user, err := h.Service.CreateUser(userName, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
