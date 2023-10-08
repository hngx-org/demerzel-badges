package models

import (
	"gorm.io/gorm"
	"time"
)

type Badge string

const (
	Beginner     Badge = "beginner"
	Intermediate Badge = "intermediate"
	Expert       Badge = "expert"
)

type SkillBadge struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SkillID   uint      `json:"skill_id"`
	Name      Badge     `json:"name"`
	MinScore  int       `json:"min_score"`
	MaxScore  int       `json:"max_score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Skill     *Skill    `json:"skill,omitempty"`
}

type UserBadge struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SkillID   uint      `json:"skill_id"`
	UserID    string    `json:"user_id" gorm:"varchar(255)"`
	BadgeID   string    `json:"badge_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User
	Skill     Skill
	Badge     SkillBadge `gorm:"foreignKey:BadgeID"`
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
