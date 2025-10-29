package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DiscordMessageLog struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	AreaID     uuid.UUID      `json:"area_id" gorm:"type:uuid;index;not null"`
	FilePath   string         `json:"file_path" gorm:"type:text"`
	SheetName  string         `json:"sheet_name" gorm:"type:text"`
	ChangeType string         `json:"change_type" gorm:"type:varchar(20)"`
	RowNumber  int            `json:"row_number"`
	Message    string         `json:"message" gorm:"type:text"`
	RowData    datatypes.JSON `json:"row_data" gorm:"type:jsonb"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
}
