package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReactions(c *gin.Context) {
	var reactions []models.Reaction
	database.DB.Preload("Service").Find(&reactions)
	c.JSON(http.StatusOK, gin.H{"data": reactions})
}

func GetReaction(c *gin.Context) {
	var reaction models.Reaction
	id := c.Param("id")

	if err := database.DB.Preload("Service").First(&reaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Réaction non trouvée"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reaction})
}

func GetServiceReactions(c *gin.Context) {
	var reactions []models.Reaction
	serviceID := c.Param("id")

	database.DB.Where("service_id = ?", serviceID).Find(&reactions)
	c.JSON(http.StatusOK, gin.H{"data": reactions})
}

func CreateReaction(c *gin.Context) {
	var input models.Reaction

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var service models.Service
	if err := database.DB.First(&service, input.ServiceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service non trouvé"})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de créer la réaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func UpdateReaction(c *gin.Context) {
	var reaction models.Reaction
	id := c.Param("id")

	if err := database.DB.First(&reaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Réaction non trouvée"})
		return
	}

	var input models.Reaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&reaction).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": reaction})
}

func DeleteReaction(c *gin.Context) {
	var reaction models.Reaction
	id := c.Param("id")

	if err := database.DB.First(&reaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Réaction non trouvée"})
		return
	}

	database.DB.Delete(&reaction)
	c.JSON(http.StatusOK, gin.H{"message": "Réaction supprimée avec succès"})
}
