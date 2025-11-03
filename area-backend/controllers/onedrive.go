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
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func OneDriveAuthStart(c *gin.Context) {
	onedriveService, err := services.NewOneDriveService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OneDrive service not configured: " + err.Error()})
		return
	}
	state := generateRandomState()

	c.SetCookie("onedrive_state", state, 300, "/", "", false, false)

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

	baseURL := strings.TrimRight(getBaseURL(), "/")
	redirectPath := "/auth/onedrive/callback"

	if code == "" {
		errorMsg := c.Query("error")
		errorDesc := c.Query("error_description")
		redirectURL := fmt.Sprintf("%s%s?error=%s&error_description=%s",
			baseURL, redirectPath, url.QueryEscape(errorMsg), url.QueryEscape(errorDesc))
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("%s%s?code=%s", baseURL, redirectPath, url.QueryEscape(code))
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

func LinkOneDriveAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
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

	tokenResp, err := onedriveService.ExchangeCodeForToken(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	userInfo, err := onedriveService.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get OneDrive user info"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	onedriveID := userInfo["id"].(string)
	onedriveEmail := ""
	if mail, ok := userInfo["mail"].(string); ok {
		onedriveEmail = mail
	} else if upn, ok := userInfo["userPrincipalName"].(string); ok {
		onedriveEmail = upn
	}

	var existingUser models.User
	if err := database.DB.Where("one_drive_id = ?", onedriveID).First(&existingUser).Error; err == nil {
		if existingUser.ID != user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "This OneDrive account is already linked to another user"})
			return
		}
	}

	user.OneDriveID = &onedriveID
	user.OneDriveEmail = &onedriveEmail
	user.OneDriveToken = &tokenResp.AccessToken
	user.OneDriveRefresh = &tokenResp.RefreshToken

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link OneDrive account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "OneDrive account linked successfully",
		"onedrive_id":    onedriveID,
		"onedrive_email": onedriveEmail,
	})
}

func UnlinkOneDriveAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.OneDriveID = nil
	user.OneDriveEmail = nil
	user.OneDriveToken = nil
	user.OneDriveRefresh = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink OneDrive account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OneDrive account unlinked successfully"})
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
