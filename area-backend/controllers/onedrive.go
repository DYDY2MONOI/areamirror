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

func OneDriveCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	_ = state

	if code == "" {
		errorMsg := c.Query("error")
		errorDesc := c.Query("error_description")
		redirectURL := fmt.Sprintf("http://localhost:3000/auth/onedrive/callback?error=%s&error_description=%s",
			errorMsg, errorDesc)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("http://localhost:3000/auth/onedrive/callback?code=%s", code)
	c.Redirect(http.StatusFound, redirectURL)
}

func OneDriveListFiles(c *gin.Context) {
	accessToken := c.GetHeader("X-OneDrive-Token")
	if accessToken == "" {
		accessToken = c.Query("access_token")
	}

	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	folderID := c.Query("folder")

	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured"})
		return
	}

	files, err := onedriveService.ListFiles(accessToken, folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    files,
	})
}

func OneDriveUploadFile(c *gin.Context) {
	accessToken := c.GetHeader("X-OneDrive-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required: " + err.Error()})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	fileName := header.Filename

	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured"})
		return
	}

	uploadResp, err := onedriveService.UploadFile(accessToken, fileName, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data":    uploadResp,
	})
}

func OneDriveDownloadFile(c *gin.Context) {
	accessToken := c.GetHeader("X-OneDrive-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	fileID := c.Param("fileId")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID required"})
		return
	}

	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured"})
		return
	}

	content, err := onedriveService.DownloadFile(accessToken, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileID))
	c.Data(http.StatusOK, "application/octet-stream", content)
}

func OneDriveDeleteFile(c *gin.Context) {
	accessToken := c.GetHeader("X-OneDrive-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	fileID := c.Param("fileId")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID required"})
		return
	}

	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured"})
		return
	}

	err = onedriveService.DeleteFile(accessToken, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

func OneDriveCreateFolder(c *gin.Context) {
	accessToken := c.GetHeader("X-OneDrive-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	var req struct {
		FolderName string `json:"folderName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured"})
		return
	}

	folder, err := onedriveService.CreateFolder(accessToken, req.FolderName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Folder created successfully",
		"data":    folder,
	})
}

func OneDriveUserInfo(c *gin.Context) {
	accessToken := c.GetHeader("X-OneDrive-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured"})
		return
	}

	userInfo, err := onedriveService.GetUserInfo(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userInfo,
	})
}
