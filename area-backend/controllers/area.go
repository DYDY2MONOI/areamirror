package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"encoding/json"
	"net/http"

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
	var area models.Area
	id := c.Param("id")

	if err := database.DB.Preload("User").Preload("Actions").Preload("Reactions").First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AREA non trouvée"})
		return
	}

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

	var user models.User
	if err := database.DB.First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user found"})
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
	var area models.Area
	id := c.Param("id")

	if err := database.DB.First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AREA non trouvée"})
		return
	}

	var input models.Area
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&area).Updates(input)

	database.DB.Preload("User").Preload("Actions").Preload("Reactions").First(&area, area.ID)

	c.JSON(http.StatusOK, gin.H{"data": area})
}

func DeleteArea(c *gin.Context) {
	var area models.Area
	id := c.Param("id")

	if err := database.DB.First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AREA non trouvée"})
		return
	}

	database.DB.Delete(&area)
	c.JSON(http.StatusOK, gin.H{"message": "AREA supprimée avec succès"})
}

func ToggleArea(c *gin.Context) {
	var area models.Area
	id := c.Param("id")

	if err := database.DB.First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AREA non trouvée"})
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
