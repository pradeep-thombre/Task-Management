package models

import "time"

type Task struct {
	ID          string    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=pending in-progress completed"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
