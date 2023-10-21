package handlers

import (
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserBadgeBySkill(c *gin.Context) {
	skillId := c.Param("skillId")
	userId := c.GetString("user_id")

	var userbadge []models.UserBadge
	var skillBadges []models.SkillBadge

	result := db.DB.Where("id = ?", skillId).Find(&skillBadges)
	if result.Error != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to get badge", map[string]interface{}{
			"error": result.Error,
		})

		return
	}

	var skillBadgeIDs []uint
	for _, badge := range skillBadges {
		skillBadgeIDs = append(skillBadgeIDs, badge.ID)
	}

	result = db.DB.Where("user_id=?", userId).
		Where("badge_id IN ?", skillBadgeIDs).
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
