package handlers

import (
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func CreateBadgeHandler(c *gin.Context) {
	type CreateBadgeRequest struct {
		SkillID  uint    `json:"skill_id"`
		Name     string  `json:"name"`
		MinScore float64 `json:"min_score"`
		MaxScore float64 `json:"max_score"`
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

	badgeName, _ := models.GetValidBadgeName(input.Name)
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

func GetBadgesForUserHandler(c *gin.Context) {

	badgeName := c.Query("badge")
	if badgeName == "" {
		badgeName = c.Query("badges")
	}

	userID := c.GetString("user_id")

	badges, err := models.GetUserBadges(db.DB, userID, badgeName)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to list badges", map[string]string{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, http.StatusOK, "User Badges", map[string]interface{}{
		"badges": badges,
	})
}

func GetUserBadgeByIDHandler(c *gin.Context) {
	badgeIDQuery := c.Param("badge_id")
	badgeID, err := strconv.ParseInt(badgeIDQuery, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid badgeID", map[string]interface{}{})
		return
	}

	userId := c.GetString("user_id")
	badge, err := models.GetUserBadgeByID(db.DB, uint(badgeID), userId)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Badge Not found", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, http.StatusOK, "User Badge", map[string]interface{}{
		"badge": badge,
	})
}

func AssignBadgeHandler(c *gin.Context) {

	type AssignBadgeReq struct {
		UserID       string `json:"user_id"`
		AssessmentID uint   `json:"assessment_id"`
	}

	type SendNewBadgeEmail struct {
		Recipient       string `json:"recipient"`
		Name            string `json:"name"`
		Skill           string `json:"skill"`
		BadgeName       string `json:"badge_name"`
		UserProfileLink string `json:"user_profile_link"`
	}

	var body AssignBadgeReq

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error()), map[string]interface{}{})
		return
	}

	isValidAssessment := models.VerifyAssessment(db.DB, body.AssessmentID)

	if !isValidAssessment {
		response.Error(c, http.StatusBadRequest, "Invalid Assessment", map[string]interface{}{
			"error": "This assessment is not valid or is under review",
		})
		return
	}

	userID := c.GetString("user_id")

	userBadge, err := models.AssignBadge(db.DB, userID, body.AssessmentID)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to assign badge", map[string]interface{}{
			"err": err.Error(),
		})

		return
	}

	emailReq := SendNewBadgeEmail{
		Recipient:       userBadge.User.Email,
		Name:            userBadge.User.FirstName,
		Skill:           userBadge.Badge.Skill.CategoryName,
		BadgeName:       string(userBadge.Badge.Name),
		UserProfileLink: "https://example.com",
	}

	client := resty.New().R()
	client.SetHeader("Content-Type", "application/json")
	client.SetBody(&emailReq)
	res, err := client.Post("https://team-titan.mrprotocoll.me/api/messaging/assessment/badge")

	if err != nil {
		response.Error(c, 500, "Something went wrong", err)
		return
	}

	if res.StatusCode() != 200 {
		response.Success(c, http.StatusCreated, "Badge Assigned Successfully, Email not Sent", map[string]interface{}{
			"badge": userBadge,
		})
		return
	}

	response.Success(c, http.StatusCreated, "Badge Assigned Successfully", map[string]interface{}{
		"badge": userBadge,
	})
}
