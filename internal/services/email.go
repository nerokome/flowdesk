package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendOTP(toEmail string, otp string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	htmlContent := fmt.Sprintf("<strong>Your Flowdesk code is: %s</strong>", otp)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{toEmail},
		Subject: "Verify your Flowdesk Account",
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	return err
}
