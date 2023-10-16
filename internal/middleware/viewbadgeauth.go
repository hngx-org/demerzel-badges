package middleware

import (
	"context"
	"demerzel-badges/pkg/response"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func CanViewBadge() gin.HandlerFunc {
	type authResponse map[string]interface{}

	type authRequest struct {
		Token      string `json:"token"`
		Permission string `json:"permission"`
	}

	return func(c *gin.Context) {
		var body authRequest
		var authResp authResponse
		token := c.GetHeader("Authorization")

		// Check Auth header was supplied
		if token == "" || len(strings.Split(token, " ")) != 2 {
			response.Error(c, http.StatusUnauthorized, "Invalid Authorization Header", map[string]interface{}{
				"Auth": "Authorization header is missing or improperly formatted",
			})
			c.Abort()
			return
		}

		body.Token = strings.Split(token, " ")[1]
		if body.Token == "" {
			response.Error(c, http.StatusUnauthorized, "Specify a bearer token", map[string]interface{}{
				"Auth": "Authorization header is missing or improperly formatted",
			})
			c.Abort()
			return
		}
		body.Permission = "badge.read"

		client := resty.New().R()
		client.SetHeader("Content-Type", "application/json")
		client.SetBody(&body)
		resp, err := client.Post("https://staging.zuri.team/api/auth/api/authorize")

		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Auth service Error", err.Error())
			c.Abort()
			return
		}
		json.Unmarshal(resp.Body(), &authResp)

		if resp.StatusCode() != 200 {
			response.Error(c, resp.StatusCode(), "You are not Authorized to access this resource", authResp["message"])
			c.Abort()
			return
		}

		user, _ := authResp["user"].(map[string]interface{})

		id, _ := user["id"].(string)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "user_id",  id))

		c.Next()
	}
}
