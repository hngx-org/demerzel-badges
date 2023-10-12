package middleware

import (
	"demerzel-badges/pkg/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"fmt"
)

const (
	jwtSecretKey = "jwt_secret"  // Replace with actual key
	UserIDKey = "userID"
)

var (
	signingMethod = jwt.SigningMethodHS256 // Replace with actual signing method
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Unauthorized access", map[string]interface{}{
				"error": "Missing Authorization header",
			})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Unauthorized access", map[string]interface{}{
				"error": "Invalid Authorization header format",
			})
			c.Abort()
			return
		}

		token := authParts[1]

		userID, err := validateAndExtractToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Unauthorized access", map[string]interface{}{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)

		c.Next()
	}
}

// validateAndExtractToken validates the JWT token and extracts user id
func validateAndExtractToken(tokenString string) (string, error) {
	
	secret := []byte(jwtSecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != signingMethod {
			return nil, jwt.ErrSignatureInvalid
		}

		return secret, nil
	})

	if err != nil {
		switch err {
		case jwt.ErrSignatureInvalid:
			return "", fmt.Errorf("Invalid token signature")
		case jwt.ErrExpired:
			return "", fmt.Errorf("Token has expired")
		default:
			return "", fmt.Errorf("Invalid or expired token")
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid user ID in token")
	}

	return userID, nil
}
