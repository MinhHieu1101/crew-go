package models

import "github.com/google/uuid"

type Roster struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	TeamID   uint      `gorm:"not null;index"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"`
	IsLeader bool      `gorm:"default:false"`
}
