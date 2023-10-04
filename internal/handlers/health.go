package handlers

import (
	"demerzel-badges/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthHandler(c *gin.Context) {
	response.Success(c, http.StatusOK, "Team Demerzel Badges Service", nil)
}
