package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Token     string         `json:"-" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null;index"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked
}
