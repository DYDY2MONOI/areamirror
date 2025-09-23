package models

type Applet struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
