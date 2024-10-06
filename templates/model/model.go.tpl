package models

import "time"

type {{.ModelName}} struct {
	ID        uint      `json:"id" gorm:"primaryKey"` // i use gorm btw
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
