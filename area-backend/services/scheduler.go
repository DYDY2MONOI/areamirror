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
	emailService    *EmailService
	discordService  *DiscordService
	weatherService  *WeatherService
	onedriveService *OneDriveService
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

	weatherService, err := NewWeatherService()
	if err != nil {
		log.Printf("Warning: Weather service not available: %v", err)
	}

	onedriveService, err := NewOneDriveService()
	if err != nil {
		log.Printf("Warning: OneDrive service not available: %v", err)
	}

	return &SchedulerService{
		emailService:    emailService,
		discordService:  discordService,
		weatherService:  weatherService,
		onedriveService: onedriveService,
	}, nil
}

func (s *SchedulerService) CheckScheduledAreas() error {
	// Check Google Calendar triggers
	if err := s.checkCalendarTriggers(); err != nil {
		log.Printf("Error checking calendar triggers: %v", err)
	}

	// Check Weather triggers
	if err := s.checkWeatherTriggers(); err != nil {
		log.Printf("Error checking weather triggers: %v", err)
	}

	// Check OneDrive triggers
	if err := s.checkOneDriveTriggers(); err != nil {
		log.Printf("Error checking OneDrive triggers: %v", err)
	}

	return nil
}

func (s *SchedulerService) checkCalendarTriggers() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Google Calendar", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch calendar areas: %v", err)
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

func (s *SchedulerService) checkWeatherTriggers() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Weather", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch weather areas: %v", err)
	}

	for _, area := range areas {
		var triggerConfig map[string]interface{}
		if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
			log.Printf("Failed to parse trigger config for area %s: %v", area.Name, err)
			continue
		}

		if s.shouldTriggerWeatherArea(area, triggerConfig) {
			if err := s.executeArea(area); err != nil {
				log.Printf("Failed to execute area %s: %v", area.Name, err)
			}
		}
	}

	return nil
}

func (s *SchedulerService) checkOneDriveTriggers() error {
	var areas []models.Area

	err := database.DB.Preload("User").Where("trigger_service = ? AND is_active = ?", "OneDrive", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch OneDrive areas: %v", err)
	}

	for _, area := range areas {
		// Skip if user not found or no OneDrive token
		if area.User.ID == 0 || area.User.OneDriveToken == nil || *area.User.OneDriveToken == "" {
			log.Printf("Skipping area %s: no user or token", area.Name)
			continue
		}

		// Check if enough time has passed since last run (avoid spam)
		if area.LastRunAt != nil {
			timeSinceLastRun := time.Since(*area.LastRunAt)
			if timeSinceLastRun < 1*time.Minute {
				continue
			}
		}

		// Determine trigger type from TriggerType field
		// For OneDrive, TriggerType will be "NewFile" or "ModifiedFile"
		triggerType := area.TriggerType
		if triggerType == "" || triggerType == "Webhook" {
			// Default for backward compatibility
			triggerType = "NewFile"
		}

		// List files from OneDrive
		filesResp, err := s.onedriveService.ListFiles(*area.User.OneDriveToken, "")
		if err != nil {
			log.Printf("Failed to list OneDrive files for area %s: %v", area.Name, err)
			continue
		}

		// Check for files based on trigger type
		for _, file := range filesResp.Value {
			shouldTrigger := false

			if triggerType == "NewFile" || triggerType == "new_file" {
				// New file: check CreatedDateTime
				if time.Since(file.CreatedDateTime) < 1*time.Minute {
					log.Printf("New file detected in OneDrive: %s (created: %v)", file.Name, file.CreatedDateTime)
					shouldTrigger = true
				}
			} else if triggerType == "ModifiedFile" || triggerType == "modified_file" {
				// Modified file: check ModifiedDateTime AND make sure it's not just created
				if time.Since(file.ModifiedDateTime) < 1*time.Minute && time.Since(file.CreatedDateTime) > 1*time.Minute {
					log.Printf("Modified file detected in OneDrive: %s (modified: %v)", file.Name, file.ModifiedDateTime)
					shouldTrigger = true
				}
			}

			if shouldTrigger {
				// Execute the action with file context
				if err := s.executeAreaWithContext(area, map[string]string{
					"fileName": file.Name,
					"fileId":   file.ID,
					"fileUrl":  file.WebURL,
				}); err != nil {
					log.Printf("Failed to execute area %s: %v", area.Name, err)
				}
				break // Only trigger once per check
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

func (s *SchedulerService) shouldTriggerWeatherArea(area models.Area, triggerConfig map[string]interface{}) bool {
	if s.weatherService == nil {
		log.Printf("Weather service not available for area %s", area.Name)
		return false
	}

	// Check if we should skip due to recent execution
	if area.LastRunAt != nil {
		timeSinceLastRun := time.Since(*area.LastRunAt)
		if timeSinceLastRun < 10*time.Minute {
			log.Printf("Weather area %s already executed recently, skipping", area.Name)
			return false
		}
	}

	// Parse weather trigger configuration
	city, ok := triggerConfig["city"].(string)
	if !ok || city == "" {
		log.Printf("City not specified for weather area %s", area.Name)
		return false
	}

	temperature, ok := triggerConfig["temperature"].(float64)
	if !ok {
		temperature = 0
	}

	condition, ok := triggerConfig["condition"].(string)
	if !ok {
		condition = ""
	}

	operator, ok := triggerConfig["operator"].(string)
	if !ok {
		operator = "greater_than"
	}

	// Create weather trigger config
	weatherConfig := WeatherTriggerConfig{
		City:        city,
		Temperature: temperature,
		Condition:   condition,
		Operator:    operator,
	}

	// Check weather trigger
	result, err := s.weatherService.CheckWeatherTrigger(weatherConfig)
	if err != nil {
		log.Printf("Failed to check weather trigger for area %s: %v", area.Name, err)
		return false
	}

	if result.Triggered {
		log.Printf("Weather trigger activated for area %s: %s", area.Name, result.Message)
		return true
	}

	return false
}

func (s *SchedulerService) executeArea(area models.Area) error {
	return s.executeAreaWithContext(area, nil)
}

func (s *SchedulerService) executeAreaWithContext(area models.Area, context map[string]string) error {
	log.Printf("Executing area: %s", area.Name)

	var actionConfig map[string]interface{}
	if err := json.Unmarshal(area.ActionConfig, &actionConfig); err != nil {
		return fmt.Errorf("failed to parse action config: %v", err)
	}

	// Replace template variables in action config
	if context != nil {
		for key, value := range actionConfig {
			if strValue, ok := value.(string); ok {
				for varName, varValue := range context {
					strValue = strings.ReplaceAll(strValue, "{{"+varName+"}}", varValue)
				}
				actionConfig[key] = strValue
			}
		}
	}

	switch area.ActionService {
	case "Gmail":
		return s.executeGmailAction(area, actionConfig)
	case "Discord":
		return s.executeDiscordAction(area, actionConfig)
	case "OneDrive":
		return s.executeOneDriveAction(area, actionConfig)
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

func (s *SchedulerService) executeOneDriveAction(area models.Area, actionConfig map[string]interface{}) error {
	if s.onedriveService == nil {
		return fmt.Errorf("OneDrive service not available")
	}

	var user models.User
	if err := database.DB.First(&user, area.UserID).Error; err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if user.OneDriveToken == nil || *user.OneDriveToken == "" {
		return fmt.Errorf("OneDrive not linked for this user")
	}

	switch area.ActionType {
	case "UploadFile", "upload":
		return s.executeOneDriveUpload(area, actionConfig, *user.OneDriveToken)
	case "CreateFolder", "createFolder":
		return s.executeOneDriveCreateFolder(area, actionConfig, *user.OneDriveToken)
	default:
		return fmt.Errorf("unsupported OneDrive action type: %s", area.ActionType)
	}
}

func (s *SchedulerService) executeOneDriveUpload(area models.Area, actionConfig map[string]interface{}, accessToken string) error {
	fileName, ok := actionConfig["fileName"].(string)
	if !ok || fileName == "" {
		fileName = fmt.Sprintf("area_%s_%s.txt", area.Name, time.Now().Format("20060102_150405"))
	}

	content, ok := actionConfig["content"].(string)
	if !ok || content == "" {
		content = fmt.Sprintf("File created by AREA: %s at %s", area.Name, time.Now().Format("2006-01-02 15:04:05"))
	}

	templateVars := map[string]string{
		"eventTitle": "Scheduled Event",
		"eventTime":  time.Now().Format("2006-01-02 15:04:05"),
		"areaName":   area.Name,
	}

	for key, value := range templateVars {
		fileName = strings.ReplaceAll(fileName, "{{"+key+"}}", value)
		content = strings.ReplaceAll(content, "{{"+key+"}}", value)
	}

	_, err := s.onedriveService.UploadFile(accessToken, fileName, []byte(content))
	if err != nil {
		return fmt.Errorf("failed to upload file to OneDrive: %v", err)
	}

	log.Printf("File uploaded to OneDrive successfully for AREA: %s", area.Name)

	area.LastRunAt = &time.Time{}
	*area.LastRunAt = time.Now()
	area.RunCount++
	area.LastRunStatus = "success"

	if err := database.DB.Save(&area).Error; err != nil {
		log.Printf("Failed to update area status: %v", err)
	}

	return nil
}

func (s *SchedulerService) executeOneDriveCreateFolder(area models.Area, actionConfig map[string]interface{}, accessToken string) error {
	folderName, ok := actionConfig["folderName"].(string)
	if !ok || folderName == "" {
		folderName = fmt.Sprintf("AREA_%s", time.Now().Format("20060102_150405"))
	}

	templateVars := map[string]string{
		"eventTitle": "Scheduled Event",
		"eventTime":  time.Now().Format("2006-01-02 15:04:05"),
		"areaName":   area.Name,
	}

	for key, value := range templateVars {
		folderName = strings.ReplaceAll(folderName, "{{"+key+"}}", value)
	}

	_, err := s.onedriveService.CreateFolder(accessToken, folderName)
	if err != nil {
		return fmt.Errorf("failed to create folder on OneDrive: %v", err)
	}

	log.Printf("Folder created on OneDrive successfully for AREA: %s", area.Name)

	area.LastRunAt = &time.Time{}
	*area.LastRunAt = time.Now()
	area.RunCount++
	area.LastRunStatus = "success"

	if err := database.DB.Save(&area).Error; err != nil {
		log.Printf("Failed to update area status: %v", err)
	}

	return nil
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
