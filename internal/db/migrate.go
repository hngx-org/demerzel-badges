package db

import (
	"demerzel-badges/internal/models"
	"os"
	"strings"
)

func Migrate() error {
	environment := os.Getenv("ENV")
	if strings.ToLower(environment) == "production" {
		return DB.AutoMigrate(
			&models.User{},
			&models.Skill{},
			&models.Assessment{},
			&models.SkillBadge{},
			&models.UserBadge{},
		)
	}

	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Permission{},
		&models.UserPermission{},
		&models.RolePermission{},
		&models.Skill{},
		&models.Assessment{},
		&models.SkillBadge{},
		&models.UserBadge{},
	)

}
