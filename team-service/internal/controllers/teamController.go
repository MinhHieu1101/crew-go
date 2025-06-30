package controllers

import (
	"net/http"
	"team-service/internal/models"
	"team-service/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string   `json:"teamName" binding:"required,min=1"`
			Managers []string `json:"managers"`
			Members  []string `json:"members"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := db.Begin()
		team := models.Team{Name: input.Name}
		if err := tx.Create(&team).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Process roles
		utils.ProcessRoster(tx, team.ID, c.GetString("userID"), true)
		utils.ProcessMembers(tx, team.ID, input.Managers, "MANAGER")
		utils.ProcessMembers(tx, team.ID, input.Members, "MEMBER")

		tx.Commit()
		c.JSON(http.StatusCreated, team)
	}
}
