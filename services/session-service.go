package services

import (
	"errors"
	"time"

	"github.com/Skyvko6607/go-api-example/config"
	"github.com/Skyvko6607/go-api-example/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const ContextUserID = "user_id"
const TokenUserID = "user_id"
const TokenJTI = "jti"

type SessionService struct {
	RedisContext *database.RedisContext
	AppSettings  *config.AppSettings
}

func (a *SessionService) IsSessionValid(c *gin.Context, userId bson.ObjectID, jti string) (bool, error) {
	val, err := a.RedisContext.Client.Get(c, GetRedisSessionKey(userId)).Result()
	if err != nil {
		return false, err
	}
	return val == jti, nil
}

func GetUserIDFromContext(c *gin.Context) (bson.ObjectID, error) {
	id, err := c.Get(ContextUserID)
	if !err {
		return bson.ObjectID{}, errors.New("user id does not exist in current session")
	}

	return id.(bson.ObjectID), nil
}

func (a *SessionService) GetValidSession(c *gin.Context) (bson.ObjectID, error) {
	session, err := c.Cookie(a.AppSettings.Auth.CookieKey)
	if err != nil {
		return bson.ObjectID{}, err
	}

	token, err := jwt.Parse(session, func(token *jwt.Token) (any, error) {
		return []byte(a.AppSettings.Auth.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return bson.ObjectID{}, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	userId, err := bson.ObjectIDFromHex(claims[TokenUserID].(string))
	if err != nil {
		return bson.ObjectID{}, err
	}

	valid, err := a.IsSessionValid(c, userId, claims[TokenJTI].(string))
	if err != nil || !valid {
		return bson.ObjectID{}, errors.New("invalid token")
	}

	return userId, nil
}

func (a *SessionService) InvalidateSession(c *gin.Context, userId bson.ObjectID) error {
	_, err := a.RedisContext.Client.Del(c, GetRedisSessionKey(userId)).Result()
	return err
}

func (a *SessionService) RenewSession(c *gin.Context, userId bson.ObjectID) (string, time.Duration, error) {
	token, jti, duration, err := a.GenerateToken(userId)
	if err != nil {
		return "", duration, err
	}

	_, err = a.RedisContext.Client.Set(c, GetRedisSessionKey(userId), jti, duration).Result()
	return token, duration, err
}

func GetRedisSessionKey(userId bson.ObjectID) string {
	return "session_" + userId.String()
}

func (a *SessionService) GenerateToken(userId bson.ObjectID) (string, string, time.Duration, error) {
	duration := time.Hour * time.Duration(a.AppSettings.Auth.ExpirationHours)
	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			TokenUserID: userId.Hex(),
			TokenJTI:    jti,
			"exp":       time.Now().Add(duration).Unix(),
		})

	tokenString, err := token.SignedString(a.AppSettings.Auth.SecretKey)
	return tokenString, jti, duration, err
}
