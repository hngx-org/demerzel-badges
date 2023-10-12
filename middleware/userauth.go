package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"demerzel-badges/pkg/response"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInUserID, exists := c.Get("loggedInUserID")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Unauthorized access", map[string]interface{}{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		requestedUserID := c.Param("user_id") 
		if loggedInUserID != requestedUserID {
			response.Error(c, http.StatusUnauthorized, "Unauthorized access", map[string]interface{}{
				"error": "You are not authorized to perform this action",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
