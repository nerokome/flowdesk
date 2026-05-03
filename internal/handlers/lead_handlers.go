package handlers

import (
	"flowdesk/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetLeadDashboard fetches the latest team stats for the main dashboard view
func GetLeadDashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Fetch latest analytics record
		var stats models.TeamAnalytics
		if err := db.Order("date desc").First(&stats).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No analytics data found"})
			return
		}

		// 2. Fetch active blockers from MemberCheckIn table (Escalated only)
		var blockers []models.MemberCheckIn
		if err := db.Where("has_blocker = ? AND is_escalated = ?", true, true).Find(&blockers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blockers"})
			return
		}

		// 3. Construct UI response
		c.JSON(http.StatusOK, gin.H{
			"summary": gin.H{
				"checked_in":    "7/9", // Static for now; count logic can be added later
				"momentum":      stats.MomentumScore,
				"blockers":      stats.ActiveBlockers,
				"tasks_shipped": 14,
			},
			"escalated_blockers": blockers,
			"momentum_scores": gin.H{
				"check_in_rate":      stats.CheckInRate,
				"blocker_resolution": stats.BlockerResolution,
				"focus_completion":   stats.FocusCompletion,
			},
		})
	}
}

// GetTeamHealth fetches a list of analytics records for trend charting
func GetTeamHealth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trends []models.TeamAnalytics
		// Fetch last 7 analytics snapshots to show the trend
		if err := db.Order("date desc").Limit(7).Find(&trends).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch health trends"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"trends": trends,
		})
	}
}


func ResolveBlocker(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id") 

		if err := db.Model(&models.MemberCheckIn{}).Where("id = ?", id).Update("has_blocker", false).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update blocker status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blocker resolved successfully"})
	}
}