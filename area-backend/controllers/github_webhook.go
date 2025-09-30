package controllers

import (
	"Golang-API-tutoriel/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GitHubWebhookController struct {
	eventProcessor *services.GitHubEventProcessor
}

func NewGitHubWebhookController() *GitHubWebhookController {
	return &GitHubWebhookController{
		eventProcessor: services.NewGitHubEventProcessor(),
	}
}

func (ghc *GitHubWebhookController) HandleWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-GitHub-Event header"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	switch eventType {
	case "push":
		if err := ghc.handlePushEvent(body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "pull_request":
		if err := ghc.handlePullRequestEvent(body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "issues":
		if err := ghc.handleIssuesEvent(body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Event type %s not handled", eventType)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (ghc *GitHubWebhookController) handlePushEvent(body []byte) error {
	var payload services.GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal push event payload: %v", err)
	}

	return ghc.eventProcessor.ProcessPushEvent(payload)
}

func (ghc *GitHubWebhookController) handlePullRequestEvent(body []byte) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal pull request event payload: %v", err)
	}

	return ghc.eventProcessor.ProcessPullRequestEvent(payload)
}

func (ghc *GitHubWebhookController) handleIssuesEvent(body []byte) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal issues event payload: %v", err)
	}

	return ghc.eventProcessor.ProcessIssuesEvent(payload)
}
