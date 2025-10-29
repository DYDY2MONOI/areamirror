package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type actionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type reactionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type serviceInfo struct {
	Name      string         `json:"name"`
	Actions   []actionInfo   `json:"actions"`
	Reactions []reactionInfo `json:"reactions"`
}

func AboutJSON(c *gin.Context) {
	clientHost := c.ClientIP()

	var services []models.Service
	if err := database.DB.Preload("Actions").Preload("Reactions").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load services"})
		return
	}

	serviceResponses := make([]serviceInfo, 0, len(services))
	for _, svc := range services {
		actions := make([]actionInfo, 0, len(svc.Actions))
		for _, act := range svc.Actions {
			actions = append(actions, actionInfo{
				Name:        act.Name,
				Description: act.Description,
			})
		}

		reactions := make([]reactionInfo, 0, len(svc.Reactions))
		for _, react := range svc.Reactions {
			reactions = append(reactions, reactionInfo{
				Name:        react.Name,
				Description: react.Description,
			})
		}

		serviceResponses = append(serviceResponses, serviceInfo{
			Name:      svc.Name,
			Actions:   actions,
			Reactions: reactions,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"client": gin.H{
			"host": clientHost,
		},
		"server": gin.H{
			"current_time": time.Now().Unix(),
			"services":     serviceResponses,
		},
	})
}
