package auth

import (
	"time"

	"github.com/Skyvko6607/go-api-learning/config"
	"github.com/Skyvko6607/go-api-learning/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Auth struct {
	RedisContext *database.RedisContext
	AppSettings  *config.AppSettings
}

func (a *Auth) IsSessionValid(c *gin.Context, userId bson.ObjectID, jti string) (bool, error) {
	val, err := a.RedisContext.Client.Get(c, GetRedisSessionKey(userId)).Result()
	if err != nil {
		return false, err
	}
	return val == jti, nil
}

func (a *Auth) InvalidateSession(c *gin.Context, userId bson.ObjectID) error {
	_, err := a.RedisContext.Client.Del(c, GetRedisSessionKey(userId)).Result()
	return err
}

func (a *Auth) RenewSession(c *gin.Context, userId bson.ObjectID) (string, time.Duration, error) {
	token, jti, duration, err := a.GenerateToken(userId)
	if err != nil {
		return "", duration, err
	}

	_, err = a.RedisContext.Client.Set(c, GetRedisSessionKey(userId), jti, duration).Result()
	return token, duration, err
}

func GetRedisSessionKey(userId bson.ObjectID) string {
	return "session_" + userId.Hex()
}

func (a *Auth) GenerateToken(userId bson.ObjectID) (string, string, time.Duration, error) {
	duration := time.Hour * time.Duration(a.AppSettings.Auth.ExpirationHours)
	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId.Hex(),
			"jti":     jti,
			"exp":     time.Now().Add(duration).Unix(),
		})

	tokenString, err := token.SignedString([]byte(a.AppSettings.Auth.SecretKey))
	return tokenString, jti, duration, err
}
