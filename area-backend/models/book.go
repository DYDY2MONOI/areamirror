package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Areas []Area `json:"areas,omitempty" gorm:"foreignKey:UserID"`
}

type Service struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description"`
	IconURL     string         `json:"icon_url"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Actions   []Action   `json:"actions,omitempty" gorm:"foreignKey:ServiceID"`
	Reactions []Reaction `json:"reactions,omitempty" gorm:"foreignKey:ServiceID"`
}

type Action struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ServiceID   uint           `json:"service_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Parameters  string         `json:"parameters" gorm:"type:json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Service Service `json:"service,omitempty"`
	Areas   []Area  `json:"areas,omitempty" gorm:"foreignKey:ActionID"`
}

type Reaction struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ServiceID   uint           `json:"service_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Parameters  string         `json:"parameters" gorm:"type:json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Service Service `json:"service,omitempty"`
	Areas   []Area  `json:"areas,omitempty" gorm:"foreignKey:ReactionID"`
}

type Area struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	ActionID    uint           `json:"action_id" gorm:"not null"`
	ReactionID  uint           `json:"reaction_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User     User     `json:"user,omitempty"`
	Action   Action   `json:"action,omitempty"`
	Reaction Reaction `json:"reaction,omitempty"`
}
