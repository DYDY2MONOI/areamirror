package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type PasswordResetEmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

type PasswordResetEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewPasswordResetEmailService() (*PasswordResetEmailService, error) {
	smtpHost := os.Getenv("PASSWORD_RESET_SMTP_HOST")
	smtpPort := os.Getenv("PASSWORD_RESET_SMTP_PORT")
	smtpUsername := os.Getenv("PASSWORD_RESET_SMTP_USERNAME")
	smtpPassword := os.Getenv("PASSWORD_RESET_SMTP_PASSWORD")
	fromEmail := os.Getenv("PASSWORD_RESET_FROM_EMAIL")

	// Valeurs par défaut pour Gmail SMTP
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	if fromEmail == "" {
		fromEmail = smtpUsername
	}

	if smtpUsername == "" || smtpPassword == "" {
		return nil, fmt.Errorf("Password reset SMTP credentials not configured. Please set PASSWORD_RESET_SMTP_USERNAME and PASSWORD_RESET_SMTP_PASSWORD environment variables")
	}

	return &PasswordResetEmailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
	}, nil
}

func (e *PasswordResetEmailService) SendPasswordResetEmail(req PasswordResetEmailRequest) error {
	// Configuration SMTP
	auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)

	// Créer le message email
	message := e.createEmailMessage(req.To, req.Subject, req.Body)

	// Envoyer l'email
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)
	
	// Connexion TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.smtpHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Authentification
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %v", err)
	}

	// Définir l'expéditeur
	if err = client.Mail(e.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// Définir le destinataire
	if err = client.Rcpt(req.To); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Envoyer le message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	log.Printf("Password reset email sent successfully to %s with subject: %s", req.To, req.Subject)
	return nil
}

func (e *PasswordResetEmailService) createEmailMessage(to, subject, body string) string {
	emailContent := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		e.fromEmail, to, subject, body,
	)

	return emailContent
}

func (e *PasswordResetEmailService) TestConnection() error {
	// Test simple de connexion SMTP
	auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.smtpHost,
	})
	if err != nil {
		return fmt.Errorf("Password reset SMTP connection failed: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.smtpHost)
	if err != nil {
		return fmt.Errorf("Password reset SMTP client creation failed: %v", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("Password reset SMTP authentication failed: %v", err)
	}

	log.Printf("Password reset SMTP connection successful. Server: %s", e.smtpHost)
	return nil
}
