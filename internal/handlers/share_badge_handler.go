package handlers

import (
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/smtp"
)

func ShareBadgeHandler(c *gin.Context) {
	type ShareBadgeRequest struct {
		UserID      string `json:"user_id"`
		BadgeID     uint   `json:"badge_id"`
		ShareMethod string `json:"share_method"` // "link" or "email"
	}

	var input ShareBadgeRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Unable to parse payload: %s", err.Error()), map[string]interface{}{})
		return
	}

	if input.ShareMethod != "link" && input.ShareMethod != "email" {
		response.Error(c, http.StatusBadRequest, "Invalid share method", map[string]interface{}{
			"error": "Share method must be 'link' or 'email'",
		})
		return
	}

	isValidBadge := models.CheckIfBadgeIsValid(db.DB, input.BadgeID)
	if !isValidBadge {
		response.Error(c, http.StatusBadRequest, "This is not a valid badge", map[string]interface{}{
			"error": "This badge does not exist or is not a valid badge",
		})
		return
	}

	userBadge, err := models.GetUserBadgeByID(db.DB, input.UserID, input.BadgeID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to retrieve user badge", map[string]interface{}{
			"err": err.Error(),
		})
		return
	}

	switch input.ShareMethod {
	case "link":
		shareableLink := generateShareableLink(userBadge)
		response.Success(c, http.StatusOK, "Shareable Link Generated Successfully", map[string]interface{}{
			"shareable_link": shareableLink,
		})

	case "email":
		err := sendBadgeByEmail(input.UserID, userBadge)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Unable to send email", map[string]interface{}{
				"err": err.Error(),
			})
			return
		}
		response.Success(c, http.StatusOK, "Badge Shared via Email Successfully", map[string]interface{}{})
	}
}

func generateShareableLink(userBadge *models.UserBadge) string {
	// Replace these placeholders with your actual domain
	domain := "domain.com"
	route := "api/share-badge/"

	// The link format: https://domain.com/api/share-badge/123?user=456
	return fmt.Sprintf("%s%s%d?user=%s", domain, route, userBadge.ID, userBadge.UserID)
}

func sendBadgeByEmail(toEmail string, badge *models.UserBadge) error {
	smtpServer := "smtp.gmail.com"
	smtpPort :=  587
	smtpUsername := "username@gmail.com"
	smtpPassword := "smtp-password"

	subject := "You've been rewarded with a badge!"
	body := fmt.Sprintf("Dear User,\n\nYou've been rewarded with a badge!\n\nBadge Details:\n%s\n\nView it here: %s", badge.String(), generateShareableLink(badge))
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	serverAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err := smtp.SendMail(serverAddr, auth, "sender_email@email.com", []string{toEmail}, []byte(message))
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
		return err
	}

	return nil
}
