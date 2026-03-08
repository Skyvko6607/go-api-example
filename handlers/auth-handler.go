package handlers

import (
	"net/http"

	"github.com/Skyvko6607/go-api-learning/auth"
	"github.com/Skyvko6607/go-api-learning/config"
	"github.com/Skyvko6607/go-api-learning/models"
	"github.com/Skyvko6607/go-api-learning/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Auth        *auth.Auth
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

	token, duration, err := h.Auth.RenewSession(c, user.ID)
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
	session, err := c.Cookie(h.AppSettings.Auth.CookieKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not authenticated"})
		return
	}

	token, err := jwt.Parse(session, func(token *jwt.Token) (any, error) {
		return []byte(h.AppSettings.Auth.SecretKey), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	userId, err := bson.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token"})
		return
	}

	valid, err := h.Auth.IsSessionValid(c, userId, claims["jti"].(string))
	if err != nil || !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token"})
		return
	}

	err = h.Auth.InvalidateSession(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to invalidate session"})
		return
	}

	c.Status(http.StatusOK)
}
