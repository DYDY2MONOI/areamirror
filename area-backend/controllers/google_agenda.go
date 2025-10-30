package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GoogleAgendaController struct {
	agendaService *services.GoogleAgendaService
}

func NewGoogleAgendaController() *GoogleAgendaController {
	agendaService, err := services.NewGoogleAgendaService()
	if err != nil {
		return &GoogleAgendaController{
			agendaService: nil,
		}
	}

	return &GoogleAgendaController{
		agendaService: agendaService,
	}
}

func (gac *GoogleAgendaController) GetAuthURL(c *gin.Context) {
	if gac.agendaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Google Agenda service not configured",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	authURL := gac.agendaService.GetAuthURL()

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"user_id":  userID,
	})
}

func (gac *GoogleAgendaController) HandleCallback(c *gin.Context) {
	if gac.agendaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Google Agenda service not configured",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code is required",
		})
		return
	}

	token, err := gac.agendaService.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to exchange code for token: " + err.Error(),
		})
		return
	}

	var existingToken models.OAuth2Token
	err = database.DB.Where("user_id = ? AND service = ?", userID, "google").First(&existingToken).Error

	if err != nil {
		oauthToken := models.OAuth2Token{
			UserID:       userID,
			Service:      "google",
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			ExpiresAt:    &token.Expiry,
		}

		if err := database.DB.Create(&oauthToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to store token: " + err.Error(),
			})
			return
		}
	} else {
		existingToken.AccessToken = token.AccessToken
		existingToken.RefreshToken = token.RefreshToken
		existingToken.TokenType = token.TokenType
		existingToken.ExpiresAt = &token.Expiry
		if err := database.DB.Save(&existingToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update token: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Google Agenda authentication successful",
		"user_id": userID,
	})
}

func (gac *GoogleAgendaController) GetUpcomingEvents(c *gin.Context) {
	if gac.agendaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Google Agenda service not configured",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	calendarID := c.DefaultQuery("calendar_id", "primary")
	maxResultsStr := c.DefaultQuery("max_results", "10")
	maxResults, err := strconv.ParseInt(maxResultsStr, 10, 64)
	if err != nil {
		maxResults = 10
	}

	events, err := gac.agendaService.GetUpcomingEvents(strconv.Itoa(int(userID)), calendarID, maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch events: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events":      events,
		"calendar_id": calendarID,
		"count":       len(events),
	})
}

func (gac *GoogleAgendaController) TestAgendaConnection(c *gin.Context) {
	if gac.agendaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Google Agenda service not configured",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	events, err := gac.agendaService.GetUpcomingEvents(strconv.Itoa(int(userID)), "primary", 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"error":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected":    true,
		"message":      "Successfully connected to Google Agenda",
		"events_count": len(events),
	})
}

func (gac *GoogleAgendaController) ListCalendars(c *gin.Context) {
	if gac.agendaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Google Agenda service not configured",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	calendars, err := gac.agendaService.ListCalendars(strconv.Itoa(int(userID)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list calendars: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"calendars": calendars,
		"count":     len(calendars),
	})
}
