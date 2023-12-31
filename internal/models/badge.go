package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Badge string

const (
	Beginner     Badge = "Beginner"
	Intermediate Badge = "Intermediate"
	Expert       Badge = "Expert"
)

func GetValidBadgeName(badgeName string) (Badge, error) {
	if strings.ToLower(string(Beginner)) == strings.ToLower(badgeName) {
		return Beginner, nil
	}

	if strings.ToLower(string(Intermediate)) == strings.ToLower(badgeName) {
		return Intermediate, nil
	}

	if strings.ToLower(string(Expert)) == strings.ToLower(badgeName) {
		return Expert, nil
	}

	return "", errors.New("invalid badge Provided")
}

type SkillBadge struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	SkillID  uint    `json:"skill_id"`
	Name     Badge   `json:"name"`
	MinScore float64 `json:"min_score"`
	MaxScore float64 `json:"max_score"`

	Skill *Skill `json:"Skill,omitempty"`
}

type SkillBadgeJson struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	SkillID  uint    `json:"skill_id"`
	Name     string  `json:"name"`
	MinScore float64 `json:"min_score"`
	MaxScore float64 `json:"max_score"`
	Skill    *Skill  `json:"Skill,omitempty"`
}

func (sB SkillBadge) MarshalJSON() ([]byte, error) {
	jsonData := SkillBadgeJson{
		ID:       sB.ID,
		SkillID:  sB.SkillID,
		Name:     strings.ToLower(string(sB.Name)),
		MinScore: sB.MinScore,
		MaxScore: sB.MaxScore,
		Skill:    sB.Skill,
	}

	return json.Marshal(jsonData)
}

func (sB SkillBadge) TableName() string {
	return "skill_badge"
}

type UserBadge struct {
	ID               uint        `json:"id" gorm:"primaryKey"`
	UserID           string      `json:"user_id" gorm:"varchar(255)"`
	BadgeID          uint        `json:"badge_id"`
	UserAssessmentID uint        `json:"user_assessment_id"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	User             *User       `json:"user,omitempty"`
	Badge            *SkillBadge `gorm:"foreignKey:BadgeID"`

	UserAssessment *UserAssessment `json:"UserAssessment"`
}

func (uB UserBadge) TableName() string {
	return "user_badge"
}

func (b Badge) IsValid() bool {
	return b == Beginner || b == Intermediate || b == Expert
}

func CreateBadge(db *gorm.DB, badge SkillBadge) (*SkillBadge, error) {
	newBadge := SkillBadge{
		SkillID:  badge.SkillID,
		Name:     badge.Name,
		MinScore: badge.MinScore,
		MaxScore: badge.MaxScore,
	}

	err := db.Create(&newBadge).Error

	return &newBadge, err
}

func BadgeExists(db *gorm.DB, skillID uint, badgeName Badge) bool {
	var existingBadge SkillBadge

	err := db.Where(&SkillBadge{SkillID: skillID, Name: badgeName}).First(&existingBadge).Error

	return err == nil
}

func AssignBadge(db *gorm.DB, userID string, assessmentID uint) (*UserBadge, error) {

	var assessmentTaken UserAssessment
	var badge SkillBadge

	err := db.Preload("Assessment").First(&assessmentTaken, assessmentID).Error

	if err != nil {
		return nil, err
	}

	err = db.Where("skill_id = ? AND ? BETWEEN min_score AND max_score", assessmentTaken.Assessment.SkillID, assessmentTaken.Score).First(&badge).Error

	if err != nil {
		return nil, err
	}

	if badge.ID == 0 {
		return nil, fmt.Errorf("badge for this assessmnt does not exist")
	}

	newUserBadge := UserBadge{
		UserID:           userID,
		BadgeID:          badge.ID,
		UserAssessmentID: assessmentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err = db.Create(&newUserBadge).Error

	if err != nil {
		return nil, err
	}

	err = db.Preload("UserAssessment").
		Preload("User").
		Preload("UserAssessment.Assessment").
		Preload("Badge").
		Preload("Badge.Skill").
		Where(&UserBadge{ID: newUserBadge.ID}).First(&newUserBadge).Error

	return &newUserBadge, err
}

func CheckIfBadgeIsValid(db *gorm.DB, badgeID uint) bool {
	var badgecheck SkillBadge
	err := db.Where(&SkillBadge{ID: badgeID}).First(&badgecheck).Error

	return err == nil
}

func VerifyAssessment(db *gorm.DB, asssessmentID uint) bool {
	var assessment_taken UserAssessment

	err := db.Where(&UserAssessment{ID: asssessmentID}).First(&assessment_taken).Error

	if assessment_taken.Status == Pending || assessment_taken.Status == Failed {
		return false
	}

	return err == nil
}

func GetUserBadgeByID(db *gorm.DB, badgeID uint, userID string) (*UserBadge, error) {
	var badge UserBadge
	result := db.Where(&UserBadge{ID: badgeID, UserID: userID}).
		Preload("User").
		Preload("Badge").
		Preload("UserAssessment.Assessment").
		Preload("Badge.Skill").
		First(&badge)

	if result.Error != nil {
		return nil, result.Error
	}

	return &badge, nil
}

func GetUserBadges(db *gorm.DB, userID string, badgeName string) ([]UserBadge, error) {
	var badges []UserBadge

	query := db
	if badgeName != "" {
		validBadgeName, err := GetValidBadgeName(badgeName)
		if err != nil {
			return nil, err
		}
		query = query.Raw("SELECT user_badge.id, user_badge.user_id, user_badge.badge_id, user_badge.user_assessment_id "+
			"FROM user_badge, skill_badge WHERE skill_badge.id = user_badge.badge_id AND skill_badge.name = ? AND user_badge.user_id = ?",
			validBadgeName, userID,
		)

	} else {
		query = query.Model(&UserBadge{}).Where(&UserBadge{UserID: userID})
	}

	result := query.Preload("UserAssessment").
		Preload("User").
		Preload("Badge").
		Preload("Badge.Skill").
		Preload("UserAssessment.Assessment").
		Find(&badges)

	if result.Error != nil {
		return nil, result.Error
	}

	return badges, nil
}
