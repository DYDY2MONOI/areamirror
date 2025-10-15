package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type SchedulerService struct {
	emailService   *EmailService
	discordService *DiscordService
	slackService   *SlackService
}

func NewSchedulerService() (*SchedulerService, error) {
	emailService, err := NewEmailService()
	if err != nil {
		log.Printf("Warning: Email service not available: %v", err)
	}

	discordService, err := NewDiscordService()
	if err != nil {
		log.Printf("Warning: Discord service not available: %v", err)
	}

	slackService, err := NewSlackService()
	if err != nil {
		log.Printf("Warning: Slack service not available: %v", err)
	}

	return &SchedulerService{
		emailService:   emailService,
		discordService: discordService,
		slackService:   slackService,
	}, nil
}

func (s *SchedulerService) CheckScheduledAreas() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Google Calendar", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch areas: %v", err)
	}

	now := time.Now()

	for _, area := range areas {
		var triggerConfig map[string]interface{}
		if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
			log.Printf("Failed to parse trigger config for area %s: %v", area.Name, err)
			continue
		}

		if s.shouldTriggerArea(area, triggerConfig, now) {
			if err := s.executeArea(area); err != nil {
				log.Printf("Failed to execute area %s: %v", area.Name, err)
			}
		}
	}

	return nil
}

func (s *SchedulerService) shouldTriggerArea(area models.Area, triggerConfig map[string]interface{}, now time.Time) bool {
	eventTimeStr, ok := triggerConfig["eventTime"].(string)
	if !ok {
		return false
	}

	eventTime, err := time.Parse(time.RFC3339, eventTimeStr)
	if err != nil {
		log.Printf("Failed to parse event time for area %s: %v", area.Name, err)
		return false
	}

	if area.LastRunAt != nil {
		timeSinceLastRun := now.Sub(*area.LastRunAt)
		if timeSinceLastRun < 5*time.Minute {
			log.Printf("Area %s already executed recently, skipping", area.Name)
			return false
		}
	}

	timeDiff := eventTime.Sub(now)
	return timeDiff >= 0 && timeDiff <= 30*time.Second
}

func (s *SchedulerService) executeArea(area models.Area) error {
	log.Printf("Executing area: %s", area.Name)

	var actionConfig map[string]interface{}
	if err := json.Unmarshal(area.ActionConfig, &actionConfig); err != nil {
		return fmt.Errorf("failed to parse action config: %v", err)
	}

	switch area.ActionService {
	case "Gmail":
		return s.executeGmailAction(area, actionConfig)
	case "Discord":
		return s.executeDiscordAction(area, actionConfig)
	case "Slack":
		return s.executeSlackAction(area, actionConfig)
	default:
		log.Printf("Unsupported action service: %s", area.ActionService)
		return nil
	}
}

func (s *SchedulerService) executeGmailAction(area models.Area, actionConfig map[string]interface{}) error {
	if s.emailService == nil {
		return fmt.Errorf("Email service not available")
	}

	toEmail, ok := actionConfig["toEmail"].(string)
	if !ok {
		return fmt.Errorf("toEmail not found in action config")
	}

	subject, ok := actionConfig["subject"].(string)
	if !ok {
		subject = "AREA Notification"
	}

	body, ok := actionConfig["body"].(string)
	if !ok {
		body = "This is an automated message from your AREA."
	}

	templateVars := map[string]string{
		"eventTitle": "Scheduled Event",
		"eventTime":  time.Now().Format("2006-01-02 15:04:05"),
		"areaName":   area.Name,
	}

	for key, value := range templateVars {
		subject = strings.ReplaceAll(subject, "{{"+key+"}}", value)
		body = strings.ReplaceAll(body, "{{"+key+"}}", value)
	}

	emailReq := EmailRequest{
		To:      toEmail,
		Subject: subject,
		Body:    body,
	}

	if err := s.emailService.SendEmail(emailReq); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email sent successfully for AREA: %s", area.Name)

	area.LastRunAt = &time.Time{}
	*area.LastRunAt = time.Now()
	area.RunCount++
	area.LastRunStatus = "success"

	if err := database.DB.Save(&area).Error; err != nil {
		log.Printf("Failed to update area status: %v", err)
	}

	log.Printf("Successfully executed area: %s", area.Name)
	return nil
}

func (s *SchedulerService) executeDiscordAction(area models.Area, actionConfig map[string]interface{}) error {
	if s.discordService == nil {
		return fmt.Errorf("Discord service not available")
	}

	webhookURL, _ := actionConfig["webhookUrl"].(string)
	if webhookURL == "" {
		if alt, ok := actionConfig["webhookURL"].(string); ok {
			webhookURL = alt
		}
	}
	if webhookURL == "" {
		return fmt.Errorf("webhookUrl not found in action config")
	}

	message, _ := actionConfig["message"].(string)
	if message == "" {
		message = fmt.Sprintf("Notification from area %s", area.Name)
	}

	templateVars := map[string]string{
		"eventTitle": "Scheduled Event",
		"eventTime":  time.Now().Format("2006-01-02 15:04:05"),
		"areaName":   area.Name,
	}

	for key, value := range templateVars {
		message = strings.ReplaceAll(message, "{{"+key+"}}", value)
	}

	if err := s.discordService.SendWebhookMessage(webhookURL, message); err != nil {
		return fmt.Errorf("failed to send discord message: %v", err)
	}

	log.Printf("Discord message sent successfully for AREA: %s", area.Name)

	area.LastRunAt = &time.Time{}
	*area.LastRunAt = time.Now()
	area.RunCount++
	area.LastRunStatus = "success"

	if err := database.DB.Save(&area).Error; err != nil {
		log.Printf("Failed to update area status: %v", err)
	}

	log.Printf("Successfully executed area: %s", area.Name)
	return nil
}

func (s *SchedulerService) executeSlackAction(area models.Area, actionConfig map[string]interface{}) error {
	if s.slackService == nil {
		return fmt.Errorf("Slack service not available")
	}

	webhookURL, _ := actionConfig["webhookUrl"].(string)
	if webhookURL == "" {
		if alt, ok := actionConfig["webhookURL"].(string); ok {
			webhookURL = alt
		}
	}
	if webhookURL == "" {
		return fmt.Errorf("webhookUrl not found in action config")
	}

	message, _ := actionConfig["message"].(string)
	if message == "" {
		message = fmt.Sprintf("📬 Notification from area: %s", area.Name)
	}

	var triggerConfig map[string]interface{}
	if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err == nil {
		templateVars := map[string]string{
			"eventTitle": getStringValue(triggerConfig, "eventTitle", "Scheduled Event"),
			"eventTime":  time.Now().Format("2006-01-02 15:04:05"),
			"areaName":   area.Name,
		}

		for key, value := range templateVars {
			message = strings.ReplaceAll(message, "{{"+key+"}}", value)
		}
	}

	messageType, _ := actionConfig["messageType"].(string)
	if messageType == "" {
		messageType = "simple"
	}

	var err error

	switch messageType {
	case "rich":
		attachment := CreateGitHubNotificationAttachment(
			area.Name,
			message,
			"AREA Scheduler",
			"",
		)
		err = s.slackService.SendRichMessage(webhookURL, "📅 Notification programmée", []Attachment{attachment})

	case "custom":
		username, _ := actionConfig["username"].(string)
		iconEmoji, _ := actionConfig["iconEmoji"].(string)

		customMsg := SlackWebhookMessage{
			Text:      message,
			Username:  username,
			IconEmoji: iconEmoji,
		}
		err = s.slackService.SendCustomMessage(webhookURL, customMsg)

	default:
		err = s.slackService.SendWebhookMessage(webhookURL, message)
	}

	if err != nil {
		return fmt.Errorf("failed to send slack message: %v", err)
	}

	log.Printf("Slack message sent successfully for AREA: %s", area.Name)

	area.LastRunAt = &time.Time{}
	*area.LastRunAt = time.Now()
	area.RunCount++
	area.LastRunStatus = "success"

	if err := database.DB.Save(&area).Error; err != nil {
		log.Printf("Failed to update area status: %v", err)
	}

	log.Printf("Successfully executed area: %s", area.Name)
	return nil
}

func getStringValue(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}

func (s *SchedulerService) StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Scheduler started - checking for scheduled areas every 30 seconds")

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			if err := s.CheckScheduledAreas(); err != nil {
				log.Printf("Error checking scheduled areas: %v", err)
			}
		}
	}
}

func (s *SchedulerService) TestScheduler(areaID string) error {
	var area models.Area
	if err := database.DB.Where("id = ?", areaID).First(&area).Error; err != nil {
		return fmt.Errorf("area not found: %v", err)
	}

	var triggerConfig map[string]interface{}
	if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
		return fmt.Errorf("failed to parse trigger config: %v", err)
	}
	triggerConfig["eventTime"] = time.Now().Format(time.RFC3339)

	updatedConfig, _ := json.Marshal(triggerConfig)
	area.TriggerConfig = updatedConfig
	database.DB.Save(&area)

	return s.executeArea(area)
}
