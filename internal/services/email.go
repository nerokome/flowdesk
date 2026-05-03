package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendOTP(toEmail string, otp string) error {
    apiKey := os.Getenv("RESEND_API_KEY")
    
    // LOG IT LOCALLY SO YOU NEVER GET STUCK
    fmt.Printf("--- DEBUG: Sending OTP %s to %s ---\n", otp, toEmail)

    client := resend.NewClient(apiKey)
    htmlContent := fmt.Sprintf("<strong>Your Flowdesk code is: %s</strong>", otp)

    params := &resend.SendEmailRequest{
        From:    "onboarding@resend.dev", // This only works for YOUR email
        To:      []string{toEmail},
        Subject: "Verify your Flowdesk Account",
        Html:    htmlContent,
    }

    _, err := client.Emails.Send(params)
    if err != nil {
        fmt.Printf("Resend Error: %v\n", err)
    
        return nil 
    }
    return nil
}
