package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type EmailService struct {
	service *gmail.Service
	config  *oauth2.Config
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
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		URL      string   `json:"url"`
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
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		URL      string   `json:"url"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"head_commit"`
}

func NewEmailService() (*EmailService, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
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

	return &EmailService{
		service: nil,
		config:  config,
	}, nil
}

func (es *EmailService) SendGitHubNotification(to, subjectTemplate, bodyTemplate string, eventData GitHubEventData) error {
	fmt.Printf("📧 Sending real email to: %s\n", to)
	fmt.Printf("📧 Repository: %s\n", eventData.Repository.FullName)
	fmt.Printf("📧 Commit: %s\n", eventData.HeadCommit.Message)
	fmt.Printf("📧 Author: %s\n", eventData.HeadCommit.Author.Name)

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
	var oauth2Token models.OAuth2Token
	if err := database.DB.Where("service = ?", "gmail").First(&oauth2Token).Error; err != nil {
		return fmt.Errorf("no Gmail OAuth2 token found: %v", err)
	}

	if oauth2Token.NeedsRefresh() {
		if err := es.refreshToken(&oauth2Token); err != nil {
			return fmt.Errorf("failed to refresh token: %v", err)
		}
	}

	token := &oauth2.Token{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		TokenType:    oauth2Token.TokenType,
		Expiry:       *oauth2Token.ExpiresAt,
	}

	client := es.config.Client(oauth2.NoContext, token)

	gmailService, err := gmail.New(client)
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %v", err)
	}

	message := es.createEmailMessage(req.To, req.Subject, req.Body)

	_, err = gmailService.Users.Messages.Send("me", &gmail.Message{
		Raw: message,
	}).Do()

	if err != nil {
		return fmt.Errorf("failed to send email via Gmail API: %v", err)
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
	return e.sendEmail(req)
}

func (es *EmailService) refreshToken(oauth2Token *models.OAuth2Token) error {
	token := &oauth2.Token{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		TokenType:    oauth2Token.TokenType,
		Expiry:       *oauth2Token.ExpiresAt,
	}

	newToken, err := es.config.TokenSource(oauth2.NoContext, token).Token()
	if err != nil {
		return fmt.Errorf("failed to refresh token: %v", err)
	}

	oauth2Token.AccessToken = newToken.AccessToken
	if newToken.RefreshToken != "" {
		oauth2Token.RefreshToken = newToken.RefreshToken
	}
	if newToken.Expiry.IsZero() {
		expiry := time.Now().Add(time.Hour)
		oauth2Token.ExpiresAt = &expiry
	} else {
		oauth2Token.ExpiresAt = &newToken.Expiry
	}

	if err := database.DB.Save(oauth2Token).Error; err != nil {
		return fmt.Errorf("failed to save refreshed token: %v", err)
	}

	return nil
}

func (es *EmailService) createEmailMessage(to, subject, body string) string {
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

func (es *EmailService) TestConnection() error {
	var oauth2Token models.OAuth2Token
	if err := database.DB.Where("service = ?", "gmail").First(&oauth2Token).Error; err != nil {
		return fmt.Errorf("no Gmail OAuth2 token found: %v", err)
	}

	if oauth2Token.NeedsRefresh() {
		if err := es.refreshToken(&oauth2Token); err != nil {
			return fmt.Errorf("failed to refresh token: %v", err)
		}
	}

	token := &oauth2.Token{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		TokenType:    oauth2Token.TokenType,
		Expiry:       *oauth2Token.ExpiresAt,
	}

	client := es.config.Client(oauth2.NoContext, token)

	gmailService, err := gmail.New(client)
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %v", err)
	}

	profile, err := gmailService.Users.GetProfile("me").Do()
	if err != nil {
		return fmt.Errorf("Gmail API connection failed: %v", err)
	}

	log.Printf("Gmail API connection successful. User: %s", profile.EmailAddress)
	return nil
}
