package email

import (
	"fmt"
	
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient struct {
	APIKey string
	From   *mail.Email
}

func NewSendGridClient(apiKey, fromEmail, fromName string) *SendGridClient {
	return &SendGridClient{
		APIKey: apiKey,
		From:   mail.NewEmail(fromName, fromEmail),
	}
}

func (c *SendGridClient) Send(toEmail, toName, subject, htmlContent string) error {
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(c.From, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(c.APIKey)
	
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("sendgrid send failed: %w", err)
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
	}
	
	return nil
}

func (c *SendGridClient) SendOTP(toEmail, code string) error {
	subject := "Your verification code"
	htmlContent := fmt.Sprintf(`<p>Your verification code is <strong>%s</strong>. It will expire in 10 minutes.</p>`, code)
	return c.Send(toEmail, "", subject, htmlContent)
}

func (c *SendGridClient) SendPasswordReset(toEmail, code string) error {
	subject := "Password reset request"
	htmlContent := fmt.Sprintf(`<p>We received a request to reset your password. Use the code below:</p><p><strong>%s</strong></p>`, code)
	return c.Send(toEmail, "", subject, htmlContent)
}