package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DiscordService struct {
	httpClient *http.Client
}

type DiscordWebhookMessage struct {
	Content string `json:"content"`
}

func NewDiscordService() (*DiscordService, error) {
	return &DiscordService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (s *DiscordService) SendWebhookMessage(webhookURL, content string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhookURL is required")
	}
	if content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	payload, err := json.Marshal(DiscordWebhookMessage{Content: content})
	if err != nil {
		return fmt.Errorf("failed to marshal discord webhook payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create discord webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}
