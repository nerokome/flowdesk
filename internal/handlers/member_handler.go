package handlers

import (
	"flowdesk/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// / SubmitCheckin handles the daily status form
func SubmitCheckin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the ID and safely handle the type
		uid, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User context missing"})
			return
		}

		var userID uuid.UUID
		var err error

		// Check if it's already a UUID or a string that needs parsing
		switch v := uid.(type) {
		case string:
			userID, err = uuid.Parse(v)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format"})
				return
			}
		case uuid.UUID:
			userID = v
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected ID type"})
			return
		}

		// 2. Bind the JSON input
		var input struct {
			ShippedYesterday   string `json:"shipped_yesterday"`
			FocusToday         string `json:"focus_today"`
			EnergyLevel        string `json:"energy_level"`
			HasBlocker         bool   `json:"has_blocker"`
			BlockerDescription string `json:"blocker_description"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// 3. Create the model
		checkIn := models.MemberCheckIn{
			UserID:             userID,
			ShippedYesterday:   input.ShippedYesterday,
			FocusToday:         input.FocusToday,
			EnergyLevel:        input.EnergyLevel,
			HasBlocker:         input.HasBlocker,
			BlockerDescription: input.BlockerDescription,
		}

		// 4. Save to Database
		if err := db.Create(&checkIn).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save check-in"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Check-in successful! 🚀"})
	}
}

func GetMemberHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User context missing"})
			return
		}

		var userID uuid.UUID
		var err error

		switch v := uid.(type) {
		case string:
			userID, err = uuid.Parse(v)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format"})
				return
			}
		case uuid.UUID:
			userID = v
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected ID type"})
			return
		}

		var history []models.MemberCheckIn
		if err := db.Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch history"})
			return
		}

		c.JSON(http.StatusOK, history)
	}
}
