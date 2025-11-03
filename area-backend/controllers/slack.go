package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func SlackAuthStart(c *gin.Context) {
	slackService := services.NewSlackOAuthService()
	config := slackService.OAuth2Config()

	state := generateSlackRandomState()
	c.SetCookie("slack_state", state, 300, "/", "", false, false)

	authURL := config.AuthCodeURL(state)

	c.JSON(http.StatusOK, gin.H{
		"authUrl": authURL,
		"state":   state,
	})
}

func SlackCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	_ = state

	baseURL := strings.TrimRight(getBaseURL(), "/")
	redirectPath := "/auth/slack/callback"

	if code == "" {
		errorMsg := c.Query("error")
		redirectURL := fmt.Sprintf("%s%s?error=%s", baseURL, redirectPath, url.QueryEscape(errorMsg))
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("%s%s?code=%s", baseURL, redirectPath, url.QueryEscape(code))
	c.Redirect(http.StatusFound, redirectURL)
}

func LinkSlackAccount(c *gin.Context) {
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

	slackService := services.NewSlackOAuthService()

	oauthResp, err := slackService.ExchangeToken(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token: " + err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	slackID := oauthResp.AuthedUser.ID
	slackTeamID := oauthResp.Team.ID

	var existingUser models.User
	if err := database.DB.Where("slack_id = ?", slackID).First(&existingUser).Error; err == nil {
		if existingUser.ID != user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "This Slack account is already linked to another user"})
			return
		}
	}

	user.SlackID = &slackID
	user.SlackTeamID = &slackTeamID
	user.SlackToken = &oauthResp.AuthedUser.AccessToken
	user.SlackBotToken = &oauthResp.AccessToken

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Slack account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Slack account linked successfully",
		"slack_id":   slackID,
		"slack_team": oauthResp.Team.Name,
	})
}

func UnlinkSlackAccount(c *gin.Context) {
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

	user.SlackID = nil
	user.SlackTeamID = nil
	user.SlackToken = nil
	user.SlackBotToken = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink Slack account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slack account unlinked successfully"})
}

func generateSlackRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
