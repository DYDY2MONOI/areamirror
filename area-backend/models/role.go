package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description"`
	Permissions string         `json:"permissions" gorm:"type:json"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Users []User `json:"users,omitempty" gorm:"many2many:user_roles;"`
}

type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	RoleID    uint      `json:"role_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

const (
	RoleAdmin  = "admin"
	RoleMember = "member"
)

const (
	PermissionReadUsers      = "users:read"
	PermissionWriteUsers     = "users:write"
	PermissionDeleteUsers    = "users:delete"
	PermissionReadAreas      = "areas:read"
	PermissionWriteAreas     = "areas:write"
	PermissionDeleteAreas    = "areas:delete"
	PermissionReadServices   = "services:read"
	PermissionWriteServices  = "services:write"
	PermissionDeleteServices = "services:delete"
	PermissionManageRoles    = "roles:manage"
	PermissionViewAnalytics  = "analytics:view"
	PermissionManageSystem   = "system:manage"
)

func GetDefaultPermissions(roleName string) []string {
	switch roleName {
	case RoleAdmin:
		return []string{
			PermissionReadUsers,
			PermissionWriteUsers,
			PermissionDeleteUsers,
			PermissionReadAreas,
			PermissionWriteAreas,
			PermissionDeleteAreas,
			PermissionReadServices,
			PermissionWriteServices,
			PermissionDeleteServices,
			PermissionManageRoles,
			PermissionViewAnalytics,
			PermissionManageSystem,
		}
	case RoleMember:
		return []string{
			PermissionReadAreas,
			PermissionWriteAreas,
			PermissionDeleteAreas,
			PermissionReadServices,
		}
	default:
		return []string{}
	}
}

