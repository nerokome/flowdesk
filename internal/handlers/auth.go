package handlers

import (
	"crypto/rand"
	"flowdesk/internal/models"
	"flowdesk/internal/services"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. Hash Password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

		// 2. Generate OTP
		otp := generateOTP()
		expiry := time.Now().Add(15 * time.Minute)

		user := models.User{
			Name:             input.Name,
			Email:            input.Email,
			Password:         string(hashedPassword),
			VerificationCode: otp,
			OTPExpiresAt:     expiry,
			IsVerified:       false,
		}

		// 3. Save to DB
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		err := services.SendOTP(user.Email, otp)
		if err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Registration successful. Please verify your email.",
			"otp":     otp,
		})
	}
}
func Login(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password required"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Please verify your email before logging in",
				"verified": false,
			})
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		})
	}
}
func VerifyOTP(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email string `json:"email" binding:"required"`
			Code  string `json:"code" binding:"required"`
		}

		// 1. Validate Input
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

		// 4. Check if OTP is expired
		if time.Now().After(user.OTPExpiresAt) {
			c.JSON(http.StatusGone, gin.H{"error": "OTP has expired. Please request a new one."})
			return
		}

		if user.VerificationCode != input.Code {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
			return
		}

		// 6. Success! Update the user in the DB
		db.Model(&user).Updates(map[string]interface{}{
			"is_verified":       true,
			"verification_code": "",
		})

		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully! You can now log in."})
	}
}
