package routes

import (
	"flowdesk/internal/handlers"
	"flowdesk/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	jwtSecret := os.Getenv("JWT_SECRET")

	// --- PUBLIC ROUTES ---
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register(db))
		auth.POST("/verify", handlers.VerifyOTP(db))
		auth.POST("/login", handlers.Login(db, jwtSecret))
	}

	// --- PROTECTED API ROUTES ---
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret)) 
	{
		// Basic Profile
		api.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Access granted to protected route"})
		})

		// Feature Modules (Separate Files)
		RegisterUploadRoutes(api, db) // Existing
		RegisterMemberRoutes(api, db) // New: For the Right-side UI
		RegisterLeadRoutes(api, db)   // New: For the Left-side UI
	}
}