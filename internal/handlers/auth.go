package handlers

import (
	"crypto/rand"
	"flowdesk/internal/auth"
	"flowdesk/internal/models"
	"flowdesk/internal/services"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func generateOTP() string {
	max := 6
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max || err != nil {
		return "123456" 
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// Register - Matches image_8.png (Step 1)
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
			Role     string `json:"role" binding:"required"` 
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill all fields correctly"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		otp := generateOTP()
		expiry := time.Now().Add(15 * time.Minute)

		user := models.User{
			Name:             input.Name,
			Email:            input.Email,
			Password:         string(hashedPassword),
			Role:             input.Role,
			VerificationCode: otp,
			OTPExpiresAt:     expiry,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "This email is already registered"})
			return
		}
		go services.SendOTP(user.Email, otp)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Step 1 complete. Verification code sent.",
			"user_id": user.ID,
		})
	}
}

// VerifyOTP - Matches image_9.png (Step 2)
func VerifyOTP(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email string `json:"email" binding:"required"`
			Code  string `json:"code" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Code are required"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if user.IsVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Account already verified"})
			return
		}

		if time.Now().After(user.OTPExpiresAt) {
			c.JSON(http.StatusGone, gin.H{"error": "OTP expired"})
			return
		}

		if user.VerificationCode != input.Code {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
			return
		}

		db.Model(&user).Updates(map[string]interface{}{
			"is_verified":       true,
			"verification_code": "",
		})

		c.JSON(http.StatusOK, gin.H{"message": "Verified! Proceed to team setup."})
	}
}

func Login(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Credentials required"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Block unverified users to respect the flow in image_9.png
		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Please verify your email before logging in",
				"verified": false,
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		tokenString, err := auth.GenerateToken(user.ID.String(), user.Role, jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		// --- JWT GENERATION END ---

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   tokenString, // Return the token to the client
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"role":  user.Role,
				"email": user.Email,
			},
		})
	}
}
