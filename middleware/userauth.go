package middleware

import (
	"demerzel-badges/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"fmt"
)

const (
	UserIDKey = "userID"
	authServiceURL = "https://auth.akuya.tech/api/authorize"
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

		userID, err := getUserID(token)
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

// getUserID sends an API request to the auth service and extracts the user ID
func getUserID(token string) (string, error) {
	
	client := resty.New()

	requestBody := map[string]string{"token": token}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(authServiceURL)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("Authorization request failed with status code: %d", resp.StatusCode())
	}

	var response struct {
		Authorized bool `json:"authorized"`
		User       struct {
			ID string `json:"id"`
		} `json:"user"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return "", err
	}

	if response.Authorized {
		return response.User.ID, nil
	}

	return "", fmt.Errorf("User is not authorized")
}
