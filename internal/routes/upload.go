package routes

import (
	"flowdesk/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUploadRoutes handles all media-related endpoints
func RegisterUploadRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	upload := rg.Group("/upload")
	{
		upload.POST("/profile-picture", handlers.UploadProfilePicture(db))

	}
}