package handlers

import (
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserBadgeHandler(c *gin.Context) {
	userId := c.Param("userId")
	skillId := c.Param("skillId")

	var userbadge []models.UserBadge

	db.DB.Where("user_id=? AND skill_id=?", userId, skillId).Find(&userbadge)
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

func GetAllUserBadgesHandler(c *gin.Context) {
	userId := c.Param("userId")

	if userId == "" {
		response.Error(c, http.StatusBadRequest,
			"error: User Id is required", map[string]interface{}{})
		return
	}

	var userBadges []models.UserBadge

	result := db.DB.Where("user_id=?", userId).
		Preload("User").
		Preload("Skill").
		Preload("Badge").
		Preload("Assessment").
		Find(&userBadges)

	if result.Error != nil {
		{
			response.Error(c, http.StatusInternalServerError, "Unable to find badges", map[string]interface{}{
				"error": result.Error,
			})
		}
	}

	if len(userBadges) == 0 {
		response.Error(c, http.StatusNotFound, "User Badges not Found", map[string]interface{}{})
		return
	}

	response.Success(c, http.StatusOK, "User Badges Retrieved Successfully", map[string]interface{}{
		"data": userBadges,
	})
}
