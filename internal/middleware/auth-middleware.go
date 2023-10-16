package middleware

import (
	"demerzel-badges/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func CanAssignBadge() gin.HandlerFunc {
	type authRequest struct {
		Token      string `json:"token"`
		Permission string `json:"permission"`
	}

	type authResponse map[string]interface{}

	return func(ctx *gin.Context) {
		var body authRequest
		var authRes authResponse
		token := ctx.GetHeader("Authorization")


// Check Auth header was supplied
		if token == "" || len(strings.Split(token, " ")) != 2 {
			response.Error(ctx, http.StatusUnauthorized, "Invalid Authorization Header", map[string]interface{}{
				"Auth": "Authorization header is missing or improperly formatted",
			})
			ctx.Abort()
			return
		}

		body.Token = strings.Split(token, " ")[1]
		if body.Token == "" {
			response.Error(ctx, http.StatusUnauthorized, "Specify a bearer token", map[string]interface{}{
				"Auth": "Authorization header is missing or improperly formatted",
			})
			ctx.Abort()
			return
		}

		body.Permission = "badge.update.own"

		client := resty.New().R()
		client.SetHeader("Content-Type", "application/json")
		client.SetBody(&body)
		res, err := client.Post("https://staging.zuri.team/api/auth/api/authorize")

		if err != nil {
			response.Error(ctx, 500, "Something went wrong", err)
			ctx.Abort()
			return
		}
		
		json.Unmarshal(res.Body(), &authRes)

		if res.StatusCode() != 200 {
			response.Error(ctx, res.StatusCode(), "You are not Authorized to access this resource", authRes["message"])
			ctx.Abort()
			return
		}
		fmt.Println("It worked")
		user, _ := authRes["user"].(map[string]interface{})

		id, _ := user["id"].(string)

		ctx.Set("user_id", id)
		ctx.Next()
	}
}
