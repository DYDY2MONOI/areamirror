package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAreaDiscordLogs(c *gin.Context) {
	areaID := c.Param("id")

	var area models.Area
	if err := database.DB.Where("id = ?", areaID).First(&area).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	if area.UserID != userID {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	limit := 50
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	var logs []models.DiscordMessageLog
	if err := database.DB.Where("area_id = ?", area.ID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}
