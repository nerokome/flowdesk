package routes

import (
	"flowdesk/internal/handlers"
	"flowdesk/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMemberRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	member := rg.Group("/member")
	member.Use(middleware.RoleBlock("member"))
	{
		member.POST("/checkin", handlers.SubmitCheckin(db))
		member.GET("/history", handlers.GetMemberHistory(db))
	}
}
