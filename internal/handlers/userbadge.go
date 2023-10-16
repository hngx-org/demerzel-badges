package handlers

import (
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserBadgeHandler(c *gin.Context) {
	skillId := c.Param("skillId")
	userId := c.GetString("user_id")


	var userbadge []models.UserBadge

	result := db.DB.Where("user_id=? AND skill_id=?", userId, skillId).
		Preload("User").
		Preload("Badge").
		Preload("Skill").
		Preload("Assessment").
		Preload("Assessment.Skill").
		Find(&userbadge)

	if result.Error != nil {
		response.Error(c, http.StatusInternalServerError, "User Badge not Found", map[string]interface{}{
			"error": result.Error,
		})
		return
	}

	if len(userbadge) == 0 {
		response.Error(c, http.StatusNotFound, "User Badge not Found", map[string]interface{}{
			"data": userbadge,
		})
		return
	}

	response.Success(c, http.StatusOK, "User Badge Retrieved Successfully", map[string]interface{}{
		"data": userbadge,
	})
}
