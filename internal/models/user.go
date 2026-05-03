package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the core account details for both Lead and Member
type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string    `json:"name"`
	AvatarURL        string    `json:"avatar_url"`
	Email            string    `gorm:"uniqueIndex" json:"email"`
	Password         string    `json:"-"`
	Role             string    `json:"role"` // "lead" or "member"
	IsVerified       bool      `gorm:"default:false" json:"is_verified"`
	VerificationCode string    `json:"-"`
	OTPExpiresAt     time.Time `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// MemberCheckIn represents the Right UI in image_11.png
// This is the "ground truth" data submitted by the team member daily.
type MemberCheckIn struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;index" json:"user_id"`

	// UI: "What did you ship yesterday?"
	ShippedYesterday string `gorm:"type:text" json:"shipped_yesterday"`

	// UI: "What are you focused on today?"
	FocusToday string `gorm:"type:text" json:"focus_today"`

	// UI: Energy level buttons (low, medium, high, peak)
	EnergyLevel string `json:"energy_level"`

	// UI: "I have a blocker" toggle + description
	HasBlocker         bool   `gorm:"default:false" json:"has_blocker"`
	BlockerDescription string `gorm:"type:text" json:"blocker_description"`
	IsEscalated        bool   `gorm:"default:false" json:"is_escalated"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TeamAnalytics represents the Left UI in image_11.png
// This stores the calculated scores the Lead sees, like "78 momentum".
type TeamAnalytics struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TeamID string    `gorm:"index" json:"team_id"`

	// UI: "team momentum score"
	MomentumScore int `json:"momentum_score"`

	// UI: Progress bars (89% rate, etc.)
	CheckInRate       int `json:"check_in_rate"`
	BlockerResolution int `json:"blocker_resolution"`
	FocusCompletion   int `json:"focus_completion"`

	// UI: "2 blockers" summary
	ActiveBlockers int `json:"active_blockers"`

	Date time.Time `json:"date"`
}

// --- HOOKS FOR AUTOMATIC UUID GENERATION ---

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (m *MemberCheckIn) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

func (t *TeamAnalytics) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
