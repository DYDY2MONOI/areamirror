package controllers

import (
	"Golang-API-tutoriel/services"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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

func (ghc *GitHubWebhookController) verifyWebhookSignature(payload []byte, signature string) bool {
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		fmt.Printf("⚠️ Warning: WEBHOOK_SECRET not configured, skipping verification\n")
		return true
	}

	if !strings.HasPrefix(signature, "sha256=") {
		fmt.Printf("❌ Invalid signature format: %s\n", signature)
		return false
	}

	signature = signature[7:]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	fmt.Printf("🔐 Expected signature: %s\n", expectedSignature)
	fmt.Printf("🔐 Received signature: %s\n", signature)

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (ghc *GitHubWebhookController) HandleWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")
	signature := c.GetHeader("X-Hub-Signature-256")

	fmt.Printf("🎣 Webhook received! Event type: %s, Signature: %s\n", eventType, signature)

	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-GitHub-Event header"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	if !ghc.verifyWebhookSignature(body, signature) {
		fmt.Printf("❌ Webhook signature verification failed\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
		return
	}

	fmt.Printf("✅ Webhook signature verified successfully\n")

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

    fmt.Printf("📦 Push payload summary: repo=%s id=%d ref=%s commits=%d head=%s\n",
        payload.Repository.FullName,
        payload.Repository.ID,
        payload.Ref,
        len(payload.Commits),
        payload.HeadCommit.ID,
    )

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

type TestPushRequest struct {
    RepositoryFullName string   `json:"repository_full_name"`
    RepositoryID       int      `json:"repository_id"`
    Ref                string   `json:"ref"`
    Commits            []struct {
        ID      string `json:"id"`
        Message string `json:"message"`
        URL     string `json:"url"`
    } `json:"commits"`
}

func TestGitHubPush(c *gin.Context) {
    var req TestPushRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    payload := services.GitHubWebhookPayload{}
    payload.Ref = req.Ref
    payload.Repository.FullName = req.RepositoryFullName
    payload.Repository.Name = req.RepositoryFullName
    payload.Repository.ID = req.RepositoryID

    for _, cm := range req.Commits {
        payload.Commits = append(payload.Commits, struct {
            ID      string   `json:"id"`
            Message string   `json:"message"`
            Author  struct {
                Name  string `json:"name"`
                Email string `json:"email"`
            } `json:"author"`
            URL      string   `json:"url"`
            Added    []string `json:"added"`
            Removed  []string `json:"removed"`
            Modified []string `json:"modified"`
        }{
            ID:      cm.ID,
            Message: cm.Message,
            URL:     cm.URL,
        })
    }

    processor := services.NewGitHubEventProcessor()
    if err := processor.ProcessPushEvent(payload); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok", "processed": true})
}
