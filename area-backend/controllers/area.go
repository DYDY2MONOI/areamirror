package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CreateAreaRequest struct {
	Name           string      `json:"name" binding:"required"`
	Description    string      `json:"description"`
	TriggerService string      `json:"triggerService" binding:"required"`
	TriggerType    string      `json:"triggerType" binding:"required"`
	ActionService  string      `json:"actionService" binding:"required"`
	ActionType     string      `json:"actionType" binding:"required"`
	TriggerConfig  interface{} `json:"triggerConfig"`
	ActionConfig   interface{} `json:"actionConfig"`
}

func GetAreas(c *gin.Context) {
	var areas []models.Area
	database.DB.Preload("User").Preload("Actions").Preload("Reactions").Find(&areas)
	c.JSON(http.StatusOK, gin.H{"data": areas})
}

func GetArea(c *gin.Context) {
	log.Println("GetArea function called")
	var area models.Area
	id := c.Param("id")
	log.Printf("Looking for area with ID: %s", id)

	if err := database.DB.First(&area, id).Error; err != nil {
		log.Printf("Area not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	log.Printf("Area found: %+v", area)
	c.JSON(http.StatusOK, gin.H{"data": area})
}

func GetUserAreas(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var areas []models.Area
	database.DB.Preload("Actions").Preload("Reactions").Where("user_id = ?", userID).Find(&areas)
	c.JSON(http.StatusOK, gin.H{"data": areas})
}

func CreateArea(c *gin.Context) {
	var req CreateAreaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	triggerConfigJSON, _ := json.Marshal(req.TriggerConfig)
	actionConfigJSON, _ := json.Marshal(req.ActionConfig)

	area := models.Area{
		UserID:         user.ID,
		Name:           req.Name,
		Description:    req.Description,
		TriggerService: req.TriggerService,
		TriggerType:    req.TriggerType,
		ActionService:  req.ActionService,
		ActionType:     req.ActionType,
		IsActive:       true,
		IsPublic:       true,
		TriggerConfig:  datatypes.JSON(triggerConfigJSON),
		ActionConfig:   datatypes.JSON(actionConfigJSON),
		TriggerIconURL: getIconUrlForService(req.TriggerService),
		ActionIconURL:  getIconUrlForService(req.ActionService),
	}

	if err := database.DB.Create(&area).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create area"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": area})
}

func getIconUrlForService(service string) string {
	switch service {
	case "Google Calendar":
		return "google-calendar.png"
	case "GitHub":
		return "github.png"
	case "Gmail":
		return "gmail.png"
	case "Discord":
		return "discord.png"
	case "Slack":
		return "slack.png"
	case "Weather":
		return "weather.png"
	case "Instagram":
		return "instagram.png"
	case "Twitter":
		return "twitter.png"
	case "YouTube":
		return "youtube.png"
	case "Spotify":
		return "spotify.png"
	case "Telegram":
		return "telegram.png"
	case "Twitch":
		return "twitch.png"
	case "Dropbox":
		return "dropbox.png"
	case "Notion":
		return "notion.png"
	default:
		return ""
	}
}

func UpdateArea(c *gin.Context) {
	log.Println("UpdateArea called with ID:", c.Param("id"))

	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("User not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var area models.Area
	id := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&area).Error; err != nil {
		log.Printf("Area not found for user %v: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	log.Printf("Area found: %+v", area)

	var req CreateAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input data: %+v", req)

	var triggerConfigJSON, actionConfigJSON datatypes.JSON
	if req.TriggerConfig != nil {
		triggerConfigBytes, _ := json.Marshal(req.TriggerConfig)
		triggerConfigJSON = datatypes.JSON(triggerConfigBytes)
	}
	if req.ActionConfig != nil {
		actionConfigBytes, _ := json.Marshal(req.ActionConfig)
		actionConfigJSON = datatypes.JSON(actionConfigBytes)
	}

	updates := map[string]interface{}{
		"name":            req.Name,
		"description":     req.Description,
		"trigger_service": req.TriggerService,
		"trigger_type":   req.TriggerType,
		"action_service": req.ActionService,
		"action_type":    req.ActionType,
	}

	if req.TriggerConfig != nil {
		updates["trigger_config"] = triggerConfigJSON
	}
	if req.ActionConfig != nil {
		updates["action_config"] = actionConfigJSON
	}

	if req.TriggerService != "" {
		updates["trigger_icon_url"] = getIconUrlForService(req.TriggerService)
	}
	if req.ActionService != "" {
		updates["action_icon_url"] = getIconUrlForService(req.ActionService)
	}

	log.Printf("Updating area with: %+v", updates)
	if err := database.DB.Model(&area).Updates(updates).Error; err != nil {
		log.Printf("Database update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update area"})
		return
	}

	database.DB.First(&area, area.ID)

	log.Printf("Updated area: %+v", area)
	c.JSON(http.StatusOK, gin.H{"data": area})
}

func DeleteArea(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var area models.Area
	id := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&area).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	database.DB.Delete(&area)
	c.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
}

func ToggleArea(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var area models.Area
	id := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&area).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	area.IsActive = !area.IsActive
	database.DB.Save(&area)

	c.JSON(http.StatusOK, gin.H{"data": area})
}

func GetPopularAreas(c *gin.Context) {
	var areas []models.Area
	database.DB.Where("is_public = ? AND is_active = ?", true, true).Limit(4).Find(&areas)

	var templates []gin.H
	for _, area := range areas {
		template := gin.H{
			"id":             area.ID,
			"title":          area.Name,
			"subtitle":       getSubtitleForArea(area),
			"description":    area.Description,
			"icon":           getIconForService(area.TriggerService),
			"gradientClass":  getGradientClassForArea(area),
			"triggerService": area.TriggerService,
			"actionService":  area.ActionService,
			"triggerIconUrl": area.TriggerIconURL,
			"actionIconUrl":  area.ActionIconURL,
			"isActive":       area.IsActive,
		}
		templates = append(templates, template)
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func GetRecommendedAreas(c *gin.Context) {
	var areas []models.Area
	database.DB.Where("is_public = ? AND is_active = ?", true, true).Offset(4).Limit(4).Find(&areas)

	var templates []gin.H
	for _, area := range areas {
		template := gin.H{
			"id":             area.ID,
			"title":          area.Name,
			"subtitle":       getSubtitleForArea(area),
			"description":    area.Description,
			"icon":           getIconForService(area.TriggerService),
			"gradientClass":  getGradientClassForArea(area),
			"triggerService": area.TriggerService,
			"actionService":  area.ActionService,
			"triggerIconUrl": area.TriggerIconURL,
			"actionIconUrl":  area.ActionIconURL,
			"isActive":       area.IsActive,
		}
		templates = append(templates, template)
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func getSubtitleForArea(area models.Area) string {
	switch area.TriggerService {
	case "Google Calendar":
		return "Calendar automation"
	case "GitHub":
		return "Development alerts"
	case "Gmail":
		return "Email automation"
	default:
		return "Automation"
	}
}

func getIconForService(service string) string {
	switch service {
	case "Google Calendar":
		return "mdi-calendar"
	case "GitHub":
		return "mdi-github"
	case "Gmail":
		return "mdi-email-outline"
	case "Discord":
		return "mdi-discord"
	case "Slack":
		return "mdi-slack"
	case "Weather":
		return "mdi-weather-partly-cloudy"
	case "Instagram":
		return "mdi-instagram"
	case "Twitter":
		return "mdi-twitter"
	case "YouTube":
		return "mdi-youtube"
	case "Spotify":
		return "mdi-music"
	case "Telegram":
		return "mdi-telegram"
	case "Twitch":
		return "mdi-twitch"
	case "Dropbox":
		return "mdi-dropbox"
	case "Notion":
		return "mdi-notebook"
	default:
		return "mdi-cog"
	}
}

func getGradientClassForArea(area models.Area) string {
	switch area.TriggerService {
	case "Google Calendar":
		return "gradient-blue"
	case "GitHub":
		return "gradient-indigo"
	case "Gmail":
		return "gradient-red"
	case "Discord":
		return "gradient-purple"
	case "Slack":
		return "gradient-green"
	case "Weather":
		return "gradient-teal"
	case "Instagram":
		return "gradient-pink"
	case "Twitter":
		return "gradient-blue"
	case "YouTube":
		return "gradient-red"
	case "Spotify":
		return "gradient-green"
	case "Telegram":
		return "gradient-blue"
	case "Twitch":
		return "gradient-purple"
	case "Dropbox":
		return "gradient-blue"
	case "Notion":
		return "gradient-gray"
	default:
		return "gradient-gray"
	}
}

func TestEmail(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailService, err := services.NewEmailService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email service not available: " + err.Error()})
		return
	}

	if err := emailService.TestConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email connection failed: " + err.Error()})
		return
	}

	emailReq := services.EmailRequest{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
	}

	if err := emailService.SendEmail(emailReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test email sent successfully!",
		"to":      req.To,
	})
}

func TestDiscord(c *gin.Context) {
	var req struct {
		WebhookURL string `json:"webhookUrl"`
		Message    string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhookURL := strings.TrimSpace(req.WebhookURL)
	if webhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhookUrl is required"})
		return
	}

	discordService, err := services.NewDiscordService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Discord service not available: " + err.Error()})
		return
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		message = "This is a test message from AREAmirror."
	}

	if err := discordService.SendWebhookMessage(webhookURL, message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send discord message: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Test Discord message sent successfully!",
		"webhookUrl": webhookURL,
	})
}

func TestScheduler(c *gin.Context) {
	areaID := c.Param("id")

	scheduler, err := services.NewSchedulerService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Scheduler not available: " + err.Error()})
		return
	}

	if err := scheduler.TestScheduler(areaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to test scheduler: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Area executed successfully",
		"area_id": areaID,
	})
}

func TestSlack(c *gin.Context) {
	var req struct {
		WebhookURL  string `json:"webhookUrl"`
		Message     string `json:"message"`
		MessageType string `json:"messageType"`
		Username    string `json:"username,omitempty"`
		Channel     string `json:"channel,omitempty"`
		IconEmoji   string `json:"iconEmoji,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhookURL := strings.TrimSpace(req.WebhookURL)
	if webhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhookUrl is required"})
		return
	}

	slackService, err := services.NewSlackService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Slack service not available: " + err.Error()})
		return
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		message = "🚀 This is a test message from AREAmirror!"
	}

	messageType := req.MessageType
	if messageType == "" {
		messageType = "simple"
	}

	var sendErr error

	switch messageType {
	case "simple":
		sendErr = slackService.SendWebhookMessage(webhookURL, message)

	case "rich":
		attachment := services.CreateGitHubNotificationAttachment(
			"test-repository",
			message,
			"AREA Bot",
			"https://github.com/test/repo/commit/abc123",
		)
		sendErr = slackService.SendRichMessage(webhookURL, "📬 Notification enrichie", []services.Attachment{attachment})

	case "custom":
		customMsg := services.SlackWebhookMessage{
			Text:      message,
			Username:  req.Username,
			Channel:   req.Channel,
			IconEmoji: req.IconEmoji,
		}
		sendErr = slackService.SendCustomMessage(webhookURL, customMsg)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messageType. Use: simple, rich, or custom"})
		return
	}

	if sendErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send slack message: " + sendErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Test Slack message sent successfully!",
		"webhookUrl":  webhookURL,
		"messageType": messageType,
	})
}

type WeatherTriggerRequest struct {
	City        string  `json:"city" binding:"required"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Operator    string  `json:"operator"`
}

type WeatherTestRequest struct {
	TriggerConfig WeatherTriggerRequest `json:"triggerConfig" binding:"required"`
}

func TestWeatherTrigger(c *gin.Context) {
	var req WeatherTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weatherService, err := services.NewWeatherService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize weather service"})
		return
	}

	config := services.WeatherTriggerConfig{
		City:        req.TriggerConfig.City,
		Temperature: req.TriggerConfig.Temperature,
		Condition:   req.TriggerConfig.Condition,
		Operator:    req.TriggerConfig.Operator,
	}

	if config.Operator == "" {
		config.Operator = "greater_than"
	}

	result, err := weatherService.TestWeatherTrigger(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"triggered": result.Triggered,
		"message":   result.Message,
		"data":      result.Data,
	})
}

func GetWeatherData(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	weatherService, err := services.NewWeatherService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize weather service"})
		return
	}

	weather, err := weatherService.GetCurrentWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    weather,
	})
}
