package repository

import (
	"team-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RosterRepository interface {
	CreateLeader(teamID uint, userID uuid.UUID) error
	AddMembers(teamID uint, userIDs []uuid.UUID) error
	RemoveMember(teamID uint, userID uuid.UUID, isManager bool) error
}

type rosterRepo struct {
	db *gorm.DB
}

func NewRosterRepo(db *gorm.DB) RosterRepository {
	return &rosterRepo{db: db}
}

func (r *rosterRepo) CreateLeader(teamID uint, userID uuid.UUID) error {
	return r.db.Create(&models.Roster{TeamID: teamID, UserID: userID, IsLeader: true}).Error
}

func (r *rosterRepo) AddMembers(teamID uint, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}
	var rs []models.Roster
	for _, id := range userIDs {
		rs = append(rs, models.Roster{TeamID: teamID, UserID: id, IsLeader: false})
	}
	return r.db.Create(&rs).Error
}

func (r *rosterRepo) RemoveMember(teamID uint, userID uuid.UUID, isManager bool) error {
	query := r.db.Where("team_id = ? AND user_id = ? AND is_leader = ?", teamID, userID, false)
	return query.Delete(&models.Roster{}).Error
}
