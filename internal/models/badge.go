package models

import "time"

type Badge string
type Status string

const (
	Beginner     Badge = "beginner"
	Intermediate Badge = "intermediate"
	Expert       Badge = "expert"
)

const (
	Pending  Status = "pending"
	Complete Status = "complete"
	Failed   Status = "failed"
)

type Skill struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CategoryName  string    `json:"category_name"`
	Description   string    `json:"description"`
	ParentSkillID uint      `json:"parent_skill_id"`
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

	Skill Skill
}

type SkillBadge struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SkillID   uint      `json:"skill_id"`
	Name      Badge     `json:"name"`
	MinScore  int       `json:"min_score"`
	MaxScore  int       `json:"max_score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Skill     Skill
}

type UserBadge struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AssessmentID uint      `json:"assessment_id"`
	UserID       string    `json:"user_id" gorm:"varchar(255)"`
	BadgeID      string    `json:"badge_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         User
	Assessment   Assessment
	Badge        SkillBadge `gorm:"foreignKey:BadgeID"`
}
