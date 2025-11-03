package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SlackService struct {
	httpClient *http.Client
}

type SlackWebhookMessage struct {
	Text        string       `json:"text,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color      string  `json:"color,omitempty"`
	Title      string  `json:"title,omitempty"`
	Text       string  `json:"text,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func NewSlackService() (*SlackService, error) {
	return &SlackService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (s *SlackService) SendWebhookMessage(webhookURL, message string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhookURL is required")
	}
	if message == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	payload := SlackWebhookMessage{
		Text: message,
	}

	return s.sendPayload(webhookURL, payload)
}

func (s *SlackService) SendRichMessage(webhookURL, message string, attachments []Attachment) error {
	if webhookURL == "" {
		return fmt.Errorf("webhookURL is required")
	}

	payload := SlackWebhookMessage{
		Text:        message,
		Attachments: attachments,
	}

	return s.sendPayload(webhookURL, payload)
}

func (s *SlackService) SendCustomMessage(webhookURL string, msg SlackWebhookMessage) error {
	if webhookURL == "" {
		return fmt.Errorf("webhookURL is required")
	}

	return s.sendPayload(webhookURL, msg)
}

func (s *SlackService) sendPayload(webhookURL string, payload SlackWebhookMessage) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack webhook payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create slack webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func CreateGitHubNotificationAttachment(repoName, commitMessage, author, commitURL string) Attachment {
	return Attachment{
		Color: "#36a64f",
		Title: fmt.Sprintf(" Nouveau commit sur %s", repoName),
		Text:  commitMessage,
		Fields: []Field{
			{
				Title: "Auteur",
				Value: author,
				Short: true,
			},
		},
		Footer:     "GitHub Notification",
		FooterIcon: "https://github.githubassets.com/favicons/favicon.png",
		Timestamp:  time.Now().Unix(),
	}
}
