package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuth2Token struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	Service      string         `json:"service" gorm:"not null;index"`
	AccessToken  string         `json:"-" gorm:"type:text;not null"`
	RefreshToken string         `json:"-" gorm:"type:text"`
	TokenType    string         `json:"token_type" gorm:"default:'Bearer'"`
	ExpiresAt    *time.Time     `json:"expires_at" gorm:"index"`
	Scope        string         `json:"scope" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (ot *OAuth2Token) IsExpired() bool {
	if ot.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ot.ExpiresAt)
}

func (ot *OAuth2Token) IsValid() bool {
	return !ot.IsExpired() && ot.AccessToken != ""
}

func (ot *OAuth2Token) NeedsRefresh() bool {
	if ot.ExpiresAt == nil {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(*ot.ExpiresAt)
}
