package middleware

import (
	"demerzel-badges/pkg/response"
	"encoding/json"
	"fmt"
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

		body.Token = strings.Split(token, " ")[1]
		body.Permission = "badge.update.own"

		client := resty.New().R()
		client.SetHeader("Content-Type", "application/json")
		client.SetBody(&body)
		res, err := client.Post("https://auth.akuya.tech/api/authorize")

		if err != nil {
			fmt.Println(err)
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
		ctx.Next()
	}
}
