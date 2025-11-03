package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GitHubIntegrationService struct {
	apiToken string
	webhookURL string
}

type GitHubWebhookConfig struct {
	Name   string   `json:"name"`
	Events []string `json:"events"`
	Config struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Secret      string `json:"secret"`
	} `json:"config"`
	Active bool `json:"active"`
}

type GitHubWebhookResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Events []string `json:"events"`
	Config struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Secret      string `json:"secret"`
	} `json:"config"`
	Active bool `json:"active"`
}

func NewGitHubIntegrationService() *GitHubIntegrationService {
	return &GitHubIntegrationService{
		apiToken: os.Getenv("GITHUB_TOKEN"),
		webhookURL: os.Getenv("WEBHOOK_URL"),
	}
}

func (gis *GitHubIntegrationService) deleteExistingWebhooks(owner, repo string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "token "+gis.apiToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var webhooks []GitHubWebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	for _, webhook := range webhooks {
		fmt.Printf(" Deleting existing webhook ID: %d, URL: %s\n", webhook.ID, webhook.Config.URL)
		if err := gis.DeleteWebhook(owner, repo, webhook.ID); err != nil {
			fmt.Printf(" Failed to delete webhook %d: %v\n", webhook.ID, err)
		} else {
			fmt.Printf(" Deleted webhook %d\n", webhook.ID)
		}
	}

	return nil
}

func (gis *GitHubIntegrationService) CreateWebhook(owner, repo string) (*GitHubWebhookResponse, error) {
	if gis.apiToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not configured")
	}
	if gis.webhookURL == "" {
		return nil, fmt.Errorf("WEBHOOK_URL not configured")
	}

	if err := gis.deleteExistingWebhooks(owner, repo); err != nil {
		fmt.Printf(" Warning: Failed to delete existing webhooks: %v\n", err)
	}

	webhookConfig := GitHubWebhookConfig{
		Name: "web",
		Events: []string{"push", "pull_request", "issues"},
		Config: struct {
			URL         string `json:"url"`
			ContentType string `json:"content_type"`
			Secret      string `json:"secret"`
		}{
			URL:         gis.webhookURL,
			ContentType: "json",
			Secret:      os.Getenv("WEBHOOK_SECRET"),
		},
		Active: true,
	}

	jsonData, err := json.Marshal(webhookConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webhook config: %v", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", owner, repo)
	fmt.Printf(" Creating webhook for %s/%s at %s\n", owner, repo, gis.webhookURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "token "+gis.apiToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf(" GitHub API response: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf(" GitHub API error: %s\n", string(body))
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var webhookResp GitHubWebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf(" Webhook created successfully! ID: %d, URL: %s\n", webhookResp.ID, webhookResp.Config.URL)
	return &webhookResp, nil
}

func (gis *GitHubIntegrationService) DeleteWebhook(owner, repo string, hookID int) error {
	if gis.apiToken == "" {
		return fmt.Errorf("GITHUB_TOKEN not configured")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks/%d", owner, repo, hookID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "token "+gis.apiToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	return nil
}

func (gis *GitHubIntegrationService) ListWebhooks(owner, repo string) ([]GitHubWebhookResponse, error) {
	if gis.apiToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not configured")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "token "+gis.apiToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var webhooks []GitHubWebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return webhooks, nil
}
