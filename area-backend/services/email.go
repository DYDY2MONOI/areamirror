package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailService struct {
	service *gmail.Service
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewEmailService() (*EmailService, error) {
	ctx := context.Background()

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	accessToken := os.Getenv("GOOGLE_ACCESS_TOKEN")
	refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || accessToken == "" {
		return nil, fmt.Errorf("Google OAuth credentials not configured")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			gmail.GmailSendScope,
			gmail.GmailReadonlyScope,
		},
		Endpoint: google.Endpoint,
	}

	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	client := config.Client(ctx, token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %v", err)
	}

	return &EmailService{
		service: service,
	}, nil
}

func (e *EmailService) SendEmail(req EmailRequest) error {
	message := e.createEmailMessage(req.To, req.Subject, req.Body)

	_, err := e.service.Users.Messages.Send("me", &gmail.Message{
		Raw: message,
	}).Do()

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email sent successfully to %s with subject: %s", req.To, req.Subject)
	return nil
}

func (e *EmailService) createEmailMessage(to, subject, body string) string {
	emailContent := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		to, subject, body,
	)

	encoded := base64.URLEncoding.EncodeToString([]byte(emailContent))

	return encoded
}

func (e *EmailService) TestConnection() error {
	profile, err := e.service.Users.GetProfile("me").Do()
	if err != nil {
		return fmt.Errorf("Gmail connection failed: %v", err)
	}

	log.Printf("Gmail connection successful. User: %s", profile.EmailAddress)
	return nil
}
