package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActions(c *gin.Context) {
	var actions []models.Action
	database.DB.Preload("Service").Find(&actions)
	c.JSON(http.StatusOK, gin.H{"data": actions})
}

func GetAction(c *gin.Context) {
	var action models.Action
	id := c.Param("id")

	if err := database.DB.Preload("Service").First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action non trouvée"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": action})
}

func GetServiceActions(c *gin.Context) {
	var actions []models.Action
	serviceID := c.Param("id")

	database.DB.Where("service_id = ?", serviceID).Find(&actions)
	c.JSON(http.StatusOK, gin.H{"data": actions})
}

func CreateAction(c *gin.Context) {
	var input models.Action

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de créer l'action"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func UpdateAction(c *gin.Context) {
	var action models.Action
	id := c.Param("id")

	if err := database.DB.First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action non trouvée"})
		return
	}

	var input models.Action
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&action).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": action})
}

func DeleteAction(c *gin.Context) {
	var action models.Action
	id := c.Param("id")

	if err := database.DB.First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action non trouvée"})
		return
	}

	database.DB.Delete(&action)
	c.JSON(http.StatusOK, gin.H{"message": "Action supprimée avec succès"})
}
