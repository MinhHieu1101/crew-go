package repository

import (
	"team-service/internal/models"

	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	FindByID(id uint) (*models.Team, error)
	Delete(id uint) error
}

type teamRepo struct {
	db *gorm.DB
}

func NewTeamRepo(db *gorm.DB) TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepo) FindByID(id uint) (*models.Team, error) {
	var t models.Team
	if err := r.db.Preload("Rosters").First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *teamRepo) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}
