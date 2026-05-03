package handlers

import (
	"flowdesk/internal/models"
	"flowdesk/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadProfilePicture(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the ID from the Middleware Context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in token"})
			return
		}

		// 2. Initialize Cloudinary
		cld := services.InitCloudinary()

		// 3. Get the file
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		// 4. Upload
		url, err := cld.UploadMedia(c.Request.Context(), src, "flowdesk/avatars")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary upload failed"})
			return
		}

		// 5. Update DB (Note: Cast userID to string if that's your DB type)
		if err := db.Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", url).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Success",
			"avatar_url": url,
		})
	}
}
