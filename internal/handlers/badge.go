package handlers

import (
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func CreateBadgeHandler(c *gin.Context) {
	type CreateBadgeRequest struct {
		SkillID  uint   `json:"skill_id"`
		Name     string `json:"name"`
		MinScore int    `json:"min_score"`
		MaxScore int    `json:"max_score"`
	}
	var input CreateBadgeRequest

	// Error if JSON request is invalid
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Unable to parse payload: %s", err.Error()), map[string]interface{}{})
		return
	}

	if input.MinScore < 0 {
		response.Error(c, http.StatusUnprocessableEntity, "Invalid input", map[string]interface{}{
			"min_score": "min_score should be at least 0",
		})

		return
	}

	if input.MinScore >= input.MaxScore {
		response.Error(c, http.StatusUnprocessableEntity, "Invalid input", map[string]interface{}{
			"max_score": "max_score should be greater than min score",
		})

		return
	}

	existingSkill, err := models.FindSkillById(db.DB, input.SkillID)
	if err != nil || existingSkill == nil {
		response.Error(c, http.StatusUnprocessableEntity, "Invalid input", map[string]interface{}{
			"skill": "no skill found matching provided ID",
		})

		return
	}

	badgeName := models.Badge(strings.ToLower(input.Name))
	if !badgeName.IsValid() {
		response.Error(c, http.StatusUnprocessableEntity, "invalid input", map[string]interface{}{
			"name": "invalid badge name",
		})

		return
	}

	badgeExists := models.BadgeExists(db.DB, input.SkillID, badgeName)
	if badgeExists {
		response.Error(c, http.StatusBadRequest, "Badge already exists", map[string]interface{}{
			"error": "Badge with name already exists for specified skill",
		})

		return
	}

	newBadge, err := models.CreateBadge(db.DB, models.SkillBadge{
		SkillID:  input.SkillID,
		Name:     badgeName,
		MinScore: input.MinScore,
		MaxScore: input.MaxScore,
	})

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to create badge", map[string]interface{}{
			"err": err.Error(),
		})

		return
	}

	response.Success(c, http.StatusCreated, "Badge Created Successfully", map[string]interface{}{
		"badge": newBadge,
	})
}

func GetUserBadgeByIDHandler(c *gin.Context) {
	badgeIDQuery := c.Param("badge_id")
	badgeID, err := strconv.ParseInt(badgeIDQuery, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid badgeID", map[string]interface{}{})
		return
	}
	badge, err := models.GetUserBadgeByID(db.DB, uint(badgeID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Badge Not found", map[string]string{})
		return
	}

	response.Success(c, http.StatusOK, "User Badge", map[string]interface{}{
		"badge": badge,
	})
	return
}

func AssignBadgeHandler(c *gin.Context) {

	type AssignBadgeReq struct {
		UserID       string `json:"user_id"`
		BadgeID      uint   `json:"badge_id"`
		AssessmentID uint   `json:"assessment_id"`
	}

	var body AssignBadgeReq

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error()), map[string]interface{}{})
		return
	}

	isValidBadge := models.CheckIfBadgeIsValid(db.DB, body.BadgeID)

	if !isValidBadge {
		response.Error(c, http.StatusBadRequest, "This is not a valid badge", map[string]interface{}{
			"error": "This badge does not exist or is not a valid badge",
		})
		return
	}

	isValidAssessment := models.VerifyAssessment(db.DB, body.AssessmentID)

	if !isValidAssessment {
		response.Error(c, http.StatusBadRequest, "Invalid Assessment", map[string]interface{}{
			"error": "This assessment is not valid or is under review",
		})
		return
	}

	userBadge, err := models.AssignBadge(db.DB, body.UserID, body.BadgeID, body.AssessmentID)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to assign badge", map[string]interface{}{
			"err": err.Error(),
		})

		return
	}

	response.Success(c, http.StatusCreated, "Badge Assigned Successfully", map[string]interface{}{
		"badge": userBadge,
	})
}
