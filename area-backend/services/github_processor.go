package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type GitHubEventProcessor struct {
	emailService *EmailService
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		ID          int    `json:"id"`
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

type AreaConfig struct {
	RepositoryID       int      `json:"repository_id"`
	RepositoryFullName string   `json:"repository_full_name"`
	NotificationTypes  []string `json:"notification_types"`
	DestinationEmail   string   `json:"destination_email"`
	SubjectTemplate    string   `json:"subject_template"`
	BodyTemplate       string   `json:"body_template"`
}

func NewGitHubEventProcessor() *GitHubEventProcessor {
	emailService, err := NewEmailService()
	if err != nil {
		log.Printf("Warning: Failed to initialize email service: %v", err)
		return &GitHubEventProcessor{
			emailService: nil,
		}
	}
	return &GitHubEventProcessor{
		emailService: emailService,
	}
}

func (gep *GitHubEventProcessor) ProcessPushEvent(payload GitHubWebhookPayload) error {
	log.Printf("Processing push event for repository: %s (ID: %d)", payload.Repository.FullName, payload.Repository.ID)

	var areas []models.Area
	err := database.DB.Where("trigger_service = ? AND is_active = ?", "github", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub areas: %v", err)
	}

	var allAreas []models.Area
	database.DB.Find(&allAreas)
	log.Printf("Total areas in database: %d", len(allAreas))
	for i, area := range allAreas {
		log.Printf("Area %d: Name=%s, TriggerService=%s, IsActive=%t", i+1, area.Name, area.TriggerService, area.IsActive)
	}

	log.Printf("Found %d GitHub areas to process", len(areas))
	for i, area := range areas {
		log.Printf("Processing area %d: %s (ID: %d)", i+1, area.Name, area.ID)
		if err := gep.processAreaForEvent(area, payload); err != nil {
			log.Printf("Error processing area %s: %v", area.Name, err)
			continue
		}
		log.Printf("Successfully processed area: %s", area.Name)
	}

	return nil
}

func (gep *GitHubEventProcessor) processAreaForEvent(area models.Area, payload GitHubWebhookPayload) error {
	var config AreaConfig
	if err := json.Unmarshal(area.TriggerConfig, &config); err != nil {
		return fmt.Errorf("failed to unmarshal trigger config: %v", err)
	}

	if config.RepositoryFullName != payload.Repository.FullName {
		log.Printf("Repository mismatch: config=%s, payload=%s", config.RepositoryFullName, payload.Repository.FullName)
		return nil
	}

	if !gep.shouldSendNotification(config.NotificationTypes, "push") {
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

	log.Printf(" Sending email to: %s", actionConfig.DestinationEmail)
	log.Printf(" Subject template: %s", subjectTemplate)
	log.Printf(" Body template: %s", bodyTemplate)

	err := gep.emailService.SendGitHubNotification(
		actionConfig.DestinationEmail,
		subjectTemplate,
		bodyTemplate,
		eventData,
	)

	if err != nil {
		log.Printf(" Failed to send email: %v", err)
		return err
	}

	log.Printf(" Email sent successfully to: %s", actionConfig.DestinationEmail)

	log.Printf("Email notification sent successfully for area %s to %s", area.Name, actionConfig.DestinationEmail)
	return nil
}

func (gep *GitHubEventProcessor) shouldSendNotification(notificationTypes []string, eventType string) bool {
	if len(notificationTypes) == 0 {
		return true
	}

	currentEvent := strings.ToLower(eventType)
	for _, notificationType := range notificationTypes {
		switch strings.ToLower(notificationType) {
		case "all":
			return true
		case currentEvent:
			return true
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
		} {
			var commits []struct {
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
			}
			for _, commit := range payload.Commits {
				commits = append(commits, struct {
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
		Pusher:     payload.Pusher,
		Ref:        payload.Ref,
		Before:     payload.Before,
		After:      payload.After,
		Created:    payload.Created,
		Deleted:    payload.Deleted,
		Forced:     payload.Forced,
		Compare:    payload.Compare,
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
