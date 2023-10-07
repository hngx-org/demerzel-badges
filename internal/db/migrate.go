package db

import "demerzel-badges/internal/models"

func Migrate() error {
	err := DB.AutoMigrate(
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
	return err
}
