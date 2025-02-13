package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email    string    `gorm:"unique" json:"email"`
	Password string    `gorm:"not null" json:"password"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}
