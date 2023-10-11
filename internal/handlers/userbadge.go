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

    //Corrected it to get the user, skill,badge and assessment that was empty
    db.DB.Where("user_id=? AND skill_id=?", userId, skillId).
        Preload("User").
        Preload("Skill").
        Preload("Badge").
        Preload("Assessment").
        Find(&userbadge)

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

