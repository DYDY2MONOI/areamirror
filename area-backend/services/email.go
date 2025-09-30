package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
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
	IsHTML  bool   `json:"is_html"`
}

type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type GitHubEventData struct {
	Repository struct {
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
	} `json:"repository"`
	Commits []struct {
		ID       string `json:"id"`
		Message  string `json:"message"`
		Author   struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		URL      string `json:"url"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Created    bool   `json:"created"`
	Deleted    bool   `json:"deleted"`
	Forced     bool   `json:"forced"`
	Compare    string `json:"compare"`
	HeadCommit struct {
		ID       string `json:"id"`
		Message  string `json:"message"`
		Author   struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		URL      string `json:"url"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"head_commit"`
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

func (es *EmailService) SendGitHubNotification(to, subjectTemplate, bodyTemplate string, eventData GitHubEventData) error {
	fmt.Printf("📧 [EMAIL SIMULÉ] Destinataire: %s\n", to)
	fmt.Printf("📧 [EMAIL SIMULÉ] Repository: %s\n", eventData.Repository.FullName)
	fmt.Printf("📧 [EMAIL SIMULÉ] Commit: %s\n", eventData.HeadCommit.Message)
	fmt.Printf("📧 [EMAIL SIMULÉ] Auteur: %s\n", eventData.HeadCommit.Author.Name)
	fmt.Printf("📧 [EMAIL SIMULÉ] Sujet: Nouveau push sur %s\n", eventData.Repository.Name)
	fmt.Printf("📧 [EMAIL SIMULÉ] Contenu: Push détecté sur le repository %s par %s\n", eventData.Repository.FullName, eventData.HeadCommit.Author.Name)
	fmt.Printf("📧 [EMAIL SIMULÉ] ==========================================\n")
	return nil

	subject, err := es.renderTemplate(subjectTemplate, eventData)
	if err != nil {
		return fmt.Errorf("failed to render subject template: %v", err)
	}

	body, err := es.renderTemplate(bodyTemplate, eventData)
	if err != nil {
		return fmt.Errorf("failed to render body template: %v", err)
	}

	emailReq := EmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}

	return es.sendEmail(emailReq)
}

func (es *EmailService) sendEmail(req EmailRequest) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	fromEmail := os.Getenv("GMAIL_USER")
	fromPassword := os.Getenv("GMAIL_PASSWORD")

	if fromEmail == "" || fromPassword == "" {
		return fmt.Errorf("GMAIL_USER and GMAIL_PASSWORD environment variables must be set")
	}

	message := fmt.Sprintf("From: %s\r\n", fromEmail)
	message += fmt.Sprintf("To: %s\r\n", req.To)
	message += fmt.Sprintf("Subject: %s\r\n", req.Subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += req.Body

	auth := smtp.PlainAuth("", fromEmail, fromPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{req.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %v", err)
	}

	return nil
}

func (es *EmailService) renderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (es *EmailService) GetDefaultPushSubjectTemplate() string {
	return "🚀 Nouvelle activité sur {{.Repository.Name}}"
}

func (es *EmailService) GetDefaultPushBodyTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>GitHub Activity Notification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f6f8fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .repo-name { font-size: 24px; font-weight: bold; color: #0366d6; margin: 0; }
        .repo-url { color: #586069; text-decoration: none; }
        .commit { background: #fff; border: 1px solid #e1e4e8; border-radius: 6px; padding: 16px; margin-bottom: 12px; }
        .commit-id { font-family: monospace; background: #f6f8fa; padding: 2px 6px; border-radius: 3px; font-size: 12px; }
        .commit-message { font-weight: bold; margin: 8px 0; }
        .commit-author { color: #586069; font-size: 14px; }
        .files { margin-top: 12px; }
        .file-list { font-size: 14px; color: #586069; }
        .added { color: #28a745; }
        .removed { color: #d73a49; }
        .modified { color: #0366d6; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #e1e4e8; color: #586069; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="repo-name">{{.Repository.Name}}</h1>
            <p><a href="{{.Repository.HTMLURL}}" class="repo-url">{{.Repository.FullName}}</a></p>
            {{if .Repository.Description}}<p>{{.Repository.Description}}</p>{{end}}
        </div>

        <h2>📝 Nouveaux commits</h2>
        {{range .Commits}}
        <div class="commit">
            <div class="commit-id">{{.ID}}</div>
            <div class="commit-message">{{.Message}}</div>
            <div class="commit-author">Par {{.Author.Name}} ({{.Author.Email}})</div>
            <div class="files">
                {{if .Added}}<div class="file-list"><span class="added">+ {{len .Added}} fichiers ajoutés</span></div>{{end}}
                {{if .Removed}}<div class="file-list"><span class="removed">- {{len .Removed}} fichiers supprimés</span></div>{{end}}
                {{if .Modified}}<div class="file-list"><span class="modified">~ {{len .Modified}} fichiers modifiés</span></div>{{end}}
            </div>
            <p><a href="{{.URL}}">Voir le commit</a></p>
        </div>
        {{end}}

        <div class="footer">
            <p>Cette notification a été envoyée automatiquement par votre AREA GitHub → Gmail.</p>
            <p>Branch: {{.Ref}} | Comparaison: <a href="{{.Compare}}">Voir les changements</a></p>
        </div>
    </div>
</body>
</html>`
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
