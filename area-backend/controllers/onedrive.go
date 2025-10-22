package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OneDriveAuthStart(c *gin.Context) {
	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured: " + err.Error()})
		return
	}
	state := generateRandomState()

	c.SetCookie("onedrive_state", state, 300, "/", "localhost", false, false)

	authURL := onedriveService.GetAuthorizationURL(state)

	c.JSON(http.StatusOK, gin.H{
		"authUrl": authURL,
		"state":   state,
	})
}
