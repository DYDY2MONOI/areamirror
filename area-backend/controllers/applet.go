package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/database"
)

func CreateApplet(c *gin.Context) {
	var applet models.Applet
	if err := c.ShouldBindJSON(&applet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	if err := database.DB.Create(&applet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création"})
		return
	}

	c.JSON(http.StatusCreated, applet)
}

func GetApplets(c *gin.Context) {
	var applets []models.Applet
	userID := c.Param("id")

	if err := database.DB.Where("user_id = ?", userID).Find(&applets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}

	c.JSON(http.StatusOK, applets)
}

func GetApplet(c *gin.Context) {
	id := c.Param("id")
	var applet models.Applet

	if err := database.DB.First(&applet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Applet non trouvé"})
		return
	}

	c.JSON(http.StatusOK, applet)
}

func UpdateApplet(c *gin.Context) {
	id := c.Param("id")
	var applet models.Applet

	if err := database.DB.First(&applet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Applet non trouvé"})
		return
	}

	if err := c.ShouldBindJSON(&applet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	if err := database.DB.Save(&applet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la modification"})
		return
	}

	c.JSON(http.StatusOK, applet)
}

func DeleteApplet(c *gin.Context) {
	id := c.Param("id")
	var applet models.Applet

	if err := database.DB.First(&applet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Applet non trouvé"})
		return
	}

	if err := database.DB.Delete(&applet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Applet supprimé"})
}
