package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAreaRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ActionID    uint   `json:"action_id" binding:"required"`
	ReactionID  uint   `json:"reaction_id" binding:"required"`
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

type CreateAreaRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ActionID    uint   `json:"action_id" binding:"required"`
	ReactionID  uint   `json:"reaction_id" binding:"required"`
}

func CreateArea(c *gin.Context) {
	var req CreateAreaRequest


	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var action models.Action
	if err := database.DB.First(&action, req.ActionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action non trouvée"})
		return
	}

	var reaction models.Reaction
	if err := database.DB.First(&reaction, req.ReactionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Réaction non trouvée"})
		return
	}

	area := models.Area{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := database.DB.Create(&area).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de créer l'AREA"})
		return
	}

	if err := database.DB.Model(&area).Association("Actions").Append(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'associer l'action"})
		return
	}

	if err := database.DB.Model(&area).Association("Reactions").Append(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'associer la réaction"})
		return
	}

	database.DB.Preload("User").Preload("Actions").Preload("Reactions").First(&area, area.ID)

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
