package handlers

import (
	"net/http"

	"github.com/Skyvko6607/go-api-example/config"
	"github.com/Skyvko6607/go-api-example/models"
	"github.com/Skyvko6607/go-api-example/repositories"
	"github.com/Skyvko6607/go-api-example/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Session     *services.SessionService
	Repo        *repositories.UserRepository
	AppSettings *config.AppSettings
}

func (h *AuthHandler) SetupEndpoints(e *gin.Engine) {
	e.POST("/auth/login", h.Login)
	e.POST("/auth/logout", h.Logout)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginDto models.LoginDTO
	err := c.Bind(&loginDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.FindByUserNameOrEmail(c, loginDto.UserNameOrEmail, loginDto.UserNameOrEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDto.Password))
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "wrong password"})
		return
	}

	token, duration, err := h.Session.RenewSession(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(h.AppSettings.Auth.CookieKey, token,
		int(duration.Seconds()),
		"/",
		h.AppSettings.Auth.Domain,
		true,
		true,
	)
	c.Status(http.StatusOK)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userId, err := services.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token"})
		return
	}

	err = h.Session.InvalidateSession(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to invalidate session"})
		return
	}

	c.Status(http.StatusOK)
}
