package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

	Phone              *string    `json:"phone" gorm:"uniqueIndex"`
	Birthday           *time.Time `json:"birthday"`
	Gender             string     `json:"gender"`
	Country            string     `json:"country"`
	Lang               string     `json:"lang" gorm:"default:'fr'"`
	PasswordResetToken *string    `json:"-" gorm:"index"`
	LoginProvider      string     `json:"login_provider" gorm:"default:'email'"`
	ProfileImage       *string    `json:"profile_image"`

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
	Areas   []Area  `json:"areas,omitempty" gorm:"many2many:area_actions;"`
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
	Areas   []Area  `json:"areas,omitempty" gorm:"many2many:area_reactions;"`
}

type Area struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:area_status;default:'enabled';not null"`
	IsPublic    bool      `json:"is_public" gorm:"not null;default:false"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`

	TriggerService string         `json:"trigger_service" gorm:"type:text;not null"`
	TriggerType    string         `json:"trigger_type" gorm:"type:text;not null"`
	TriggerConfig  datatypes.JSON `json:"trigger_config" gorm:"type:jsonb;not null;default:'{}'"`

	Conditions datatypes.JSON `json:"conditions" gorm:"type:jsonb;not null;default:'[]'"`

	ActionService string         `json:"action_service" gorm:"type:text;not null"`
	ActionType    string         `json:"action_type" gorm:"type:text;not null"`
	ActionConfig  datatypes.JSON `json:"action_config" gorm:"type:jsonb;not null;default:'{}'"`

	ScheduleCron    string `json:"schedule_cron" gorm:"type:text"`
	RateLimitPerMin *int   `json:"rate_limit_per_min" gorm:"check:rate_limit_per_min IS NULL OR rate_limit_per_min > 0"`
	DedupWindowSec  *int   `json:"dedup_window_sec" gorm:"check:dedup_window_sec IS NULL OR dedup_window_sec >= 0"`
	RetryMax        int    `json:"retry_max" gorm:"not null;default:3;check:retry_max >= 0"`
	RetryBackoffMs  int    `json:"retry_backoff_ms" gorm:"not null;default:1000;check:retry_backoff_ms >= 0"`

	LastRunStatus string     `json:"last_run_status" gorm:"type:run_status;default:'idle';not null"`
	LastRunAt     *time.Time `json:"last_run_at" gorm:"type:timestamptz"`
	NextRunAt     *time.Time `json:"next_run_at" gorm:"type:timestamptz"`
	RunCount      int64      `json:"run_count" gorm:"not null;default:0"`
	LastError     string     `json:"last_error" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:now()"`

	DedupKeyTemplate string `json:"dedup_key_template" gorm:"type:text"`

	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Actions   []Action   `json:"actions,omitempty" gorm:"many2many:area_actions;"`
	Reactions []Reaction `json:"reactions,omitempty" gorm:"many2many:area_reactions;"`
}
