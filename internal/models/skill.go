package models

import (
	"gorm.io/gorm"
	"time"
)

type Status string

const (
	Pending  Status = "pending"
	Complete Status = "complete"
	Failed   Status = "failed"
)

type Skill struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CategoryName  string    `json:"category_name"`
	Description   string    `json:"description"`
	ParentSkillID *uint     `json:"parent_skill_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Assessment struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	SkillID         uint      `json:"skill_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	DurationMinutes uint      `json:"duration_minutes"`
	PassScore       uint      `json:"pass_score"`
	Status          Status    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Skill           *Skill    `json:"skill,omitempty"`
}

type UserAssessment struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID         	string    `json:"user_id"`
	AssessmentID    uint      `json:"assessment_id"`
	Score       	uint      `json:"score"`
	TimeSpent 		uint      `json:"time_spent"`
	SubmissionDate  time.Time `json:"submisssion_date"`
	Status          Status    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Assessment Assessment
}

func FindSkillById(db *gorm.DB, skillID uint) (*Skill, error) {
	var existingSkill Skill
	err := db.Model(&Skill{}).First(&existingSkill, skillID).Error

	if err != nil {
		return nil, err
	}

	return &existingSkill, nil
}
