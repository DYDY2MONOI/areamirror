package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type GmailTokenRequest struct {
	Code string `json:"code" binding:"required"`
}

type GmailTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SetupGmailOAuth2(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			gmail.GmailSendScope,
			gmail.GmailReadonlyScope,
		},
		Endpoint: google.Endpoint,
	}

	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"message":  "Please visit the auth_url to authorize Gmail access",
	})
}

func StoreGmailToken(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GmailTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			gmail.GmailSendScope,
			gmail.GmailReadonlyScope,
		},
		Endpoint: google.Endpoint,
	}

	token, err := config.Exchange(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	var oauth2Token models.OAuth2Token
	err = database.DB.Where("user_id = ? AND service = ?", userID, "gmail").First(&oauth2Token).Error

	expiry := token.Expiry
	if expiry.IsZero() {
		expiry = time.Now().Add(time.Hour)
	}

	if err != nil {
		oauth2Token = models.OAuth2Token{
			UserID:       userID.(uint),
			Service:      "gmail",
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			ExpiresAt:    &expiry,
			Scope:        "https://www.googleapis.com/auth/gmail.send https://www.googleapis.com/auth/gmail.readonly",
		}

		if err := database.DB.Create(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
			return
		}
	} else {
		oauth2Token.AccessToken = token.AccessToken
		if token.RefreshToken != "" {
			oauth2Token.RefreshToken = token.RefreshToken
		}
		oauth2Token.ExpiresAt = &expiry

		if err := database.DB.Save(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token"})
			return
		}
	}

	c.JSON(http.StatusOK, GmailTokenResponse{
		Success: true,
		Message: "Gmail OAuth2 token stored successfully",
	})
}

func TestGmailConnection(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var oauth2Token models.OAuth2Token
	if err := database.DB.Where("user_id = ? AND service = ?", userID, "gmail").First(&oauth2Token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gmail OAuth2 token not found. Please set up Gmail OAuth2 first."})
		return
	}

	emailService, err := services.NewEmailService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email service"})
		return
	}

	if err := emailService.TestConnection(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Gmail connection test failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Gmail connection test successful",
	})
}

func GetGmailTokenStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var oauth2Token models.OAuth2Token
	if err := database.DB.Where("user_id = ? AND service = ?", userID, "gmail").First(&oauth2Token).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"has_token": false,
			"message":   "No Gmail OAuth2 token found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_token":     true,
		"is_valid":      oauth2Token.IsValid(),
		"needs_refresh": oauth2Token.NeedsRefresh(),
		"expires_at":    oauth2Token.ExpiresAt,
		"created_at":    oauth2Token.CreatedAt,
	})
}

func RevokeGmailToken(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := database.DB.Where("user_id = ? AND service = ?", userID, "gmail").Delete(&models.OAuth2Token{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Gmail OAuth2 token revoked successfully",
	})
}
