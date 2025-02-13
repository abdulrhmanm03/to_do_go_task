package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Status      bool      `gorm:"default:false" json:"status"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
}

func (task *Task) BeforeCreate(tx *gorm.DB) error {
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	return nil
}
