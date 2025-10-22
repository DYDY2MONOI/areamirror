package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarService struct {
	config *oauth2.Config
}

type CalendarEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	CalendarID  string    `json:"calendarId"`
}

type CalendarTriggerConfig struct {
	CalendarID    string `json:"calendarId"`
	EventTitle    string `json:"eventTitle"`
	TimeBefore    string `json:"timeBefore"`
	CheckInterval string `json:"checkInterval"`
}

func NewGoogleCalendarService() (*GoogleCalendarService, error) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			calendar.CalendarReadonlyScope,
			calendar.CalendarScope,
		},
		Endpoint: google.Endpoint,
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("Google Calendar service not configured: missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}

	return &GoogleCalendarService{
		config: config,
	}, nil
}

func (s *GoogleCalendarService) GetAuthURL() string {
	return s.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (s *GoogleCalendarService) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	return s.config.Exchange(context.Background(), code)
}

func (s *GoogleCalendarService) GetClient(token *oauth2.Token) (*calendar.Service, error) {
	client := s.config.Client(context.Background(), token)
	return calendar.NewService(context.Background(), option.WithHTTPClient(client))
}

func (s *GoogleCalendarService) GetUpcomingEvents(userID string, calendarID string, maxResults int64) ([]CalendarEvent, error) {
	token, err := s.getUserToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token: %v", err)
	}

	client, err := s.GetClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar client: %v", err)
	}

	now := time.Now().Format(time.RFC3339)
	timeMax := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	events, err := client.Events.List(calendarID).
		TimeMin(now).
		TimeMax(timeMax).
		MaxResults(maxResults).
		SingleEvents(true).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to list events: %v", err)
	}

	var calendarEvents []CalendarEvent
	for _, event := range events.Items {
		startTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			startTime, err = time.Parse("2006-01-02", event.Start.Date)
			if err != nil {
				log.Printf("Failed to parse event start time: %v", err)
				continue
			}
		}

		endTime, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			endTime, err = time.Parse("2006-01-02", event.End.Date)
			if err != nil {
				log.Printf("Failed to parse event end time: %v", err)
				continue
			}
		}

		calendarEvents = append(calendarEvents, CalendarEvent{
			ID:          event.Id,
			Title:       event.Summary,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: event.Description,
			Location:    event.Location,
			CalendarID:  calendarID,
		})
	}

	return calendarEvents, nil
}

func (s *GoogleCalendarService) CheckForUpcomingEvents(userID string, config CalendarTriggerConfig) ([]CalendarEvent, error) {
	events, err := s.GetUpcomingEvents(userID, config.CalendarID, 10)
	if err != nil {
		return nil, err
	}

	var matchingEvents []CalendarEvent
	now := time.Now()

	for _, event := range events {
		// Check if event title matches filter (if specified)
		if config.EventTitle != "" && event.Title != config.EventTitle {
			continue
		}

		// Check if event is starting within the specified time window
		timeBefore, err := parseDuration(config.TimeBefore)
		if err != nil {
			log.Printf("Failed to parse timeBefore duration: %v", err)
			continue
		}

		timeUntilEvent := event.StartTime.Sub(now)
		if timeUntilEvent > 0 && timeUntilEvent <= timeBefore {
			matchingEvents = append(matchingEvents, event)
		}
	}

	return matchingEvents, nil
}

func (s *GoogleCalendarService) getUserToken(userID string) (*oauth2.Token, error) {
	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	var oauthToken models.OAuth2Token
	err = database.DB.Where("user_id = ? AND service = ?", userIDInt, "google_calendar").First(&oauthToken).Error
	if err != nil {
		return nil, fmt.Errorf("user %s needs to authenticate with Google Calendar", userID)
	}

	token := &oauth2.Token{
		AccessToken:  oauthToken.AccessToken,
		RefreshToken: oauthToken.RefreshToken,
		TokenType:    oauthToken.TokenType,
		Expiry:       *oauthToken.ExpiresAt,
	}

	return token, nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	if durationStr == "" {
		return 15 * time.Minute, nil
	}

	switch durationStr {
	case "5m":
		return 5 * time.Minute, nil
	case "15m":
		return 15 * time.Minute, nil
	case "30m":
		return 30 * time.Minute, nil
	case "1h":
		return 1 * time.Hour, nil
	case "2h":
		return 2 * time.Hour, nil
	case "1d":
		return 24 * time.Hour, nil
	default:
		return time.ParseDuration(durationStr)
	}
}

