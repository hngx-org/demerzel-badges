package models

import "time"

type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	SectionOrder string    `json:"section_order"`
	Password     string    `json:"password"`
	ProfilePic   string    `json:"profile_pic"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u User) TableName() string {
	return "user"
}

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (r Role) TableName() string {
	return "role"
}

type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoleID    uint      `json:"role_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User
	Role Role
}

type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p Permission) TableName() string {
	return "permission"
}

type UserPermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id"`
	PermissionId uint      `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User       User
	Permission Permission
}

func (uP UserPermission) TableName() string {
	return "user_permission"
}

type RolePermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoleID       uint      `json:"role_id"`
	PermissionId uint      `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Role       Role
	Permission Permission
}

func (rP RolePermission) TableName() string {
	return "roles_permissions"
}
