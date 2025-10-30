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

type GoogleAgendaService struct {
	config *oauth2.Config
}

type AgendaEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	CalendarID  string    `json:"calendarId"`
	Attendees   []string  `json:"attendees"`
	Status      string    `json:"status"`
}

type AgendaTriggerConfig struct {
	CalendarID    string `json:"calendarId"`
	EventTitle    string `json:"eventTitle"`
	TimeBefore    string `json:"timeBefore"`
	CheckInterval string `json:"checkInterval"`
	EventStatus   string `json:"eventStatus"`
}

func NewGoogleAgendaService() (*GoogleAgendaService, error) {
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
		return nil, fmt.Errorf("Google Agenda service not configured: missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}

	return &GoogleAgendaService{
		config: config,
	}, nil
}

func (s *GoogleAgendaService) GetAuthURL() string {
	return s.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (s *GoogleAgendaService) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	return s.config.Exchange(context.Background(), code)
}

func (s *GoogleAgendaService) GetClient(token *oauth2.Token) (*calendar.Service, error) {
	client := s.config.Client(context.Background(), token)
	return calendar.NewService(context.Background(), option.WithHTTPClient(client))
}

func (s *GoogleAgendaService) GetUpcomingEvents(userID string, calendarID string, maxResults int64) ([]AgendaEvent, error) {
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

	log.Printf("Attempting to list events for calendar ID: %s", calendarID)

	events, err := client.Events.List(calendarID).
		TimeMin(now).
		TimeMax(timeMax).
		MaxResults(maxResults).
		SingleEvents(true).
		OrderBy("startTime").
		Do()

	if err != nil {
		log.Printf("Calendar API error for calendar ID '%s': %v", calendarID, err)

		if calendarID != "primary" {
			log.Printf("Trying fallback to 'primary' calendar")
			events, err = client.Events.List("primary").
				TimeMin(now).
				TimeMax(timeMax).
				MaxResults(maxResults).
				SingleEvents(true).
				OrderBy("startTime").
				Do()

			if err != nil {
				log.Printf("Calendar API error for 'primary' calendar: %v", err)
				return nil, fmt.Errorf("failed to list events for both '%s' and 'primary' calendars: %v", calendarID, err)
			}
		} else {
			return nil, fmt.Errorf("failed to list events for calendar '%s': %v", calendarID, err)
		}
	}

	var agendaEvents []AgendaEvent
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

		var attendees []string
		for _, attendee := range event.Attendees {
			attendees = append(attendees, attendee.Email)
		}

		agendaEvents = append(agendaEvents, AgendaEvent{
			ID:          event.Id,
			Title:       event.Summary,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: event.Description,
			Location:    event.Location,
			CalendarID:  calendarID,
			Attendees:   attendees,
			Status:      event.Status,
		})
	}

	return agendaEvents, nil
}

func (s *GoogleAgendaService) CheckForUpcomingEvents(userID string, config AgendaTriggerConfig) ([]AgendaEvent, error) {
	events, err := s.GetUpcomingEvents(userID, config.CalendarID, 10)
	if err != nil {
		return nil, err
	}

	var matchingEvents []AgendaEvent
	now := time.Now()

	for _, event := range events {
		if config.EventTitle != "" && event.Title != config.EventTitle {
			continue
		}

		if config.EventStatus != "" && event.Status != config.EventStatus {
			continue
		}

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

func (s *GoogleAgendaService) getUserToken(userID string) (*oauth2.Token, error) {
	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	var oauthToken models.OAuth2Token
	err = database.DB.Where("user_id = ? AND service = ?", userIDInt, "google").First(&oauthToken).Error
	if err != nil {
		return nil, fmt.Errorf("user %s needs to authenticate with Google", userID)
	}

	token := &oauth2.Token{
		AccessToken:  oauthToken.AccessToken,
		RefreshToken: oauthToken.RefreshToken,
		TokenType:    oauthToken.TokenType,
		Expiry:       *oauthToken.ExpiresAt,
	}

	return token, nil
}

func (s *GoogleAgendaService) ListCalendars(userID string) ([]string, error) {
	token, err := s.getUserToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token: %v", err)
	}

	client, err := s.GetClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar client: %v", err)
	}

	calendarList, err := client.CalendarList.List().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list calendars: %v", err)
	}

	var calendarIDs []string
	for _, calendar := range calendarList.Items {
		calendarIDs = append(calendarIDs, calendar.Id)
		log.Printf("Available calendar: %s (%s)", calendar.Id, calendar.Summary)
	}

	return calendarIDs, nil
}
