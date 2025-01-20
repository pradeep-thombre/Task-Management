package dbmodels

import "time"

type TaskSchema struct {
	ID          string    `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"not null" json:"description"`
	Status      string    `gorm:"not null" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
