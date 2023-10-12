package db

import (
	"demerzel-badges/internal/models"
	"os"
	"strings"
)

func Migrate() error {
	environment := os.Getenv("ENV")
	if strings.ToLower(environment) == "production" {
		return nil
	}

	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserPermission{},
		&models.RolePermission{},
		&models.Skill{},
		&models.Assessment{},
		&models.UserAssessment{},
		&models.SkillBadge{},
		&models.UserBadge{},
	)

}
