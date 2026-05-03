package routes

import (
	"flowdesk/internal/handlers"
	"flowdesk/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterLeadRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	lead := rg.Group("/lead")
	lead.Use(middleware.RoleBlock("lead"))
	{
		lead.GET("/dashboard", handlers.GetLeadDashboard(db))
		lead.GET("/analytics/trends", handlers.GetTeamHealth(db))
		lead.POST("/blockers/resolve", handlers.ResolveBlocker(db))
	}
}
