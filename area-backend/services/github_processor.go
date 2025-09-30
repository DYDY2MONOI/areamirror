package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"encoding/json"
	"fmt"
	"log"
)

type GitHubEventProcessor struct {
	emailService *EmailService
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
		Description string `json:"description"`
	} `json:"repository"`
	Commits []struct {
		ID       string `json:"id"`
		Message  string `json:"message"`
		Author   struct {
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
		URL      string   `json:"url"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"head_commit"`
}

type AreaConfig struct {
	RepositoryID      int      `json:"repository_id"`
	NotificationTypes  []string `json:"notification_types"`
	DestinationEmail   string   `json:"destination_email"`
	SubjectTemplate    string   `json:"subject_template"`
	BodyTemplate       string   `json:"body_template"`
}

func NewGitHubEventProcessor() *GitHubEventProcessor {
	return &GitHubEventProcessor{
		emailService: NewEmailService(),
	}
}

func (gep *GitHubEventProcessor) ProcessPushEvent(payload GitHubWebhookPayload) error {
	log.Printf("Processing push event for repository: %s (ID: %d)", payload.Repository.FullName, payload.Repository.ID)

	var areas []models.Area
	err := database.DB.Where("trigger_service = ? AND is_active = ?", "github", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub areas: %v", err)
	}

	for _, area := range areas {
		if err := gep.processAreaForEvent(area, payload); err != nil {
			log.Printf("Error processing area %s: %v", area.Name, err)
			continue
		}
	}

	return nil
}

func (gep *GitHubEventProcessor) processAreaForEvent(area models.Area, payload GitHubWebhookPayload) error {
	var config AreaConfig
	if err := json.Unmarshal(area.TriggerConfig, &config); err != nil {
		return fmt.Errorf("failed to unmarshal trigger config: %v", err)
	}

	if config.RepositoryID != payload.Repository.ID {
		return nil
	}

	if !gep.shouldSendNotification(config.NotificationTypes, payload) {
		return nil
	}

	var actionConfig AreaConfig
	if err := json.Unmarshal(area.ActionConfig, &actionConfig); err != nil {
		return fmt.Errorf("failed to unmarshal action config: %v", err)
	}

	eventData := gep.convertToEventData(payload)

	subjectTemplate := actionConfig.SubjectTemplate
	if subjectTemplate == "" {
		subjectTemplate = gep.emailService.GetDefaultPushSubjectTemplate()
	}

	bodyTemplate := actionConfig.BodyTemplate
	if bodyTemplate == "" {
		bodyTemplate = gep.emailService.GetDefaultPushBodyTemplate()
	}

	err := gep.emailService.SendGitHubNotification(
		actionConfig.DestinationEmail,
		subjectTemplate,
		bodyTemplate,
		eventData,
	)

	if err != nil {
		return fmt.Errorf("failed to send email notification: %v", err)
	}

	log.Printf("Email notification sent successfully for area %s to %s", area.Name, actionConfig.DestinationEmail)
	return nil
}

func (gep *GitHubEventProcessor) shouldSendNotification(notificationTypes []string, payload GitHubWebhookPayload) bool {
	if len(notificationTypes) == 0 {
		return true
	}

	for _, notificationType := range notificationTypes {
		switch notificationType {
		case "push":
			return true
		case "pull_request":
			return false
		case "issues":
			return false
		default:
			continue
		}
	}

	return false
}

func (gep *GitHubEventProcessor) convertToEventData(payload GitHubWebhookPayload) GitHubEventData {
	return GitHubEventData{
		Repository: struct {
			Name        string `json:"name"`
			FullName    string `json:"full_name"`
			HTMLURL     string `json:"html_url"`
			Description string `json:"description"`
		}{
			Name:        payload.Repository.Name,
			FullName:    payload.Repository.FullName,
			HTMLURL:     payload.Repository.HTMLURL,
			Description: payload.Repository.Description,
		},
		Commits: func() []struct {
			ID       string `json:"id"`
			Message  string `json:"message"`
			Author   struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
			URL      string   `json:"url"`
			Added    []string `json:"added"`
			Removed  []string `json:"removed"`
			Modified []string `json:"modified"`
		} {
			var commits []struct {
				ID       string `json:"id"`
				Message  string `json:"message"`
				Author   struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
				URL      string   `json:"url"`
				Added    []string `json:"added"`
				Removed  []string `json:"removed"`
				Modified []string `json:"modified"`
			}
			for _, commit := range payload.Commits {
				commits = append(commits, struct {
					ID       string `json:"id"`
					Message  string `json:"message"`
					Author   struct {
						Name  string `json:"name"`
						Email string `json:"email"`
					} `json:"author"`
					URL      string   `json:"url"`
					Added    []string `json:"added"`
					Removed  []string `json:"removed"`
					Modified []string `json:"modified"`
				}{
					ID:       commit.ID,
					Message:  commit.Message,
					Author:   commit.Author,
					URL:      commit.URL,
					Added:    commit.Added,
					Removed:  commit.Removed,
					Modified: commit.Modified,
				})
			}
			return commits
		}(),
		Pusher: payload.Pusher,
		Ref:    payload.Ref,
		Before: payload.Before,
		After:  payload.After,
		Created: payload.Created,
		Deleted: payload.Deleted,
		Forced:  payload.Forced,
		Compare: payload.Compare,
		HeadCommit: payload.HeadCommit,
	}
}

func (gep *GitHubEventProcessor) ProcessPullRequestEvent(payload map[string]interface{}) error {
	log.Printf("Processing pull request event")
	return nil
}

func (gep *GitHubEventProcessor) ProcessIssuesEvent(payload map[string]interface{}) error {
	log.Printf("Processing issues event")
	return nil
}
