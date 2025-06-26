package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"type:ENUM('MANAGER','MEMBER');not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
