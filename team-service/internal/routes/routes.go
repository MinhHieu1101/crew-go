package routes

import (
	"team-service/internal/controllers"
	"team-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/teams")
	api.Use(middleware.Authenticate(db))

	api.POST("/", middleware.Authorize("MANAGER"), controllers.CreateTeam(db))
	api.GET("/:teamId", controllers.GetTeam(db))
	api.DELETE("/:teamId", middleware.Authorize("MANAGER"), controllers.RemoveTeam(db))
	api.POST("/:teamId/members", middleware.Authorize("MANAGER"), controllers.AddMember(db))
	api.DELETE("/:teamId/members/:memberId", middleware.Authorize("MANAGER"), controllers.RemoveMember(db))
	api.POST("/:teamId/managers", middleware.Authorize("MANAGER"), controllers.AddManager(db))
	api.DELETE("/:teamId/managers/:managerId", middleware.Authorize("MANAGER"), controllers.RemoveManager(db))
}
