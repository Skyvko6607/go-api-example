package middleware

import (
	"net/http"
	"strings"

	"github.com/Skyvko6607/go-api-example/services"
	"github.com/gin-gonic/gin"
)

func GetAuthMiddleware(a *services.SessionService) gin.HandlerFunc {
	authEndpoints := []string{
		"/auth/logout",
	}

	return func(c *gin.Context) {
		for _, endpoint := range authEndpoints {
			if strings.HasPrefix(c.Request.URL.Path, endpoint) {
				userId, err := a.GetValidSession(c)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": err.Error(),
					})
					return
				}

				c.Set(services.ContextUserID, userId)
				c.Next()
				return
			}
		}

		c.Next()
	}
}
