package models

import "time"

type Team struct {
	ID        uint     `gorm:"primaryKey;autoIncrement"`
	Name      string   `gorm:"size:255;uniqueIndex;not null"`
	Rosters   []Roster `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
