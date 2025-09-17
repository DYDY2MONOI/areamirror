package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetServices(c *gin.Context) {
	var services []models.Service
	database.DB.Preload("Actions").Preload("Reactions").Find(&services)
	c.JSON(http.StatusOK, gin.H{"data": services})
}

func GetService(c *gin.Context) {
	var service models.Service
	id := c.Param("id")

	if err := database.DB.Preload("Actions").Preload("Reactions").First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": service})
}

func CreateService(c *gin.Context) {
	var input models.Service

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de créer le service"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func UpdateService(c *gin.Context) {
	var service models.Service
	id := c.Param("id")

	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}

	var input models.Service
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&service).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": service})
}

func DeleteService(c *gin.Context) {
	var service models.Service
	id := c.Param("id")

	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}

	database.DB.Delete(&service)
	c.JSON(http.StatusOK, gin.H{"message": "Service supprimé avec succès"})
}