package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CreateAreaRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	TriggerService string `json:"triggerService" binding:"required"`
	TriggerType    string `json:"triggerType" binding:"required"`
	ActionService  string `json:"actionService" binding:"required"`
	ActionType     string `json:"actionType" binding:"required"`
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
	var areas []models.Area
	userID := c.Param("id")

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
		TriggerConfig:  datatypes.JSON(`{}`),
		ActionConfig:   datatypes.JSON(`{}`),
	}

	if err := database.DB.Create(&area).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create area"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": area})
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
