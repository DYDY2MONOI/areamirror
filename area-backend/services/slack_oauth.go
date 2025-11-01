package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type SlackOAuthService struct{}

type SlackOAuthResponse struct {
	OK            bool   `json:"ok"`
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	Scope         string `json:"scope"`
	BotUserID     string `json:"bot_user_id"`
	AppID         string `json:"app_id"`
	Team          SlackTeam `json:"team"`
	Enterprise    interface{} `json:"enterprise"`
	AuthedUser    SlackAuthedUser `json:"authed_user"`
	Error         string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type SlackTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackAuthedUser struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type SlackUserInfo struct {
	OK    bool `json:"ok"`
	User  SlackUserProfile `json:"user"`
	Error string `json:"error"`
}

type SlackUserProfile struct {
	ID       string `json:"id"`
	TeamID   string `json:"team_id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Profile  struct {
		Email string `json:"email"`
	} `json:"profile"`
}

type SlackMessageRequest struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackMessageResponse struct {
	OK      bool   `json:"ok"`
	Channel string `json:"channel"`
	TS      string `json:"ts"`
	Error   string `json:"error"`
}

func NewSlackOAuthService() *SlackOAuthService {
	return &SlackOAuthService{}
}

func (s *SlackOAuthService) OAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("SLACK_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SLACK_REDIRECT_URI"),
		Scopes: []string{
			"channels:read",
			"channels:history",
			"chat:write",
			"users:read",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/v2/authorize",
			TokenURL: "https://slack.com/api/oauth.v2.access",
		},
	}
}

func (s *SlackOAuthService) ExchangeToken(ctx context.Context, code string) (*SlackOAuthResponse, error) {
	config := s.OAuth2Config()

	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		config.ClientID,
		config.ClientSecret,
		code,
		config.RedirectURL,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", config.Endpoint.TokenURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var oauthResp SlackOAuthResponse
	if err := json.Unmarshal(body, &oauthResp); err != nil {
		return nil, err
	}

	if !oauthResp.OK {
		return nil, errors.New(fmt.Sprintf("slack oauth error: %s - %s", oauthResp.Error, oauthResp.ErrorDescription))
	}

	return &oauthResp, nil
}

func (s *SlackOAuthService) GetUserInfo(ctx context.Context, botToken, userID string) (*SlackUserInfo, error) {
	url := fmt.Sprintf("https://slack.com/api/users.info?user=%s", userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+botToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo SlackUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	if !userInfo.OK {
		return nil, errors.New(fmt.Sprintf("slack api error: %s", userInfo.Error))
	}

	return &userInfo, nil
}

func (s *SlackOAuthService) SendMessage(botToken, channel, message string) error {
	messageReq := SlackMessageRequest{
		Channel: channel,
		Text:    message,
	}

	jsonData, err := json.Marshal(messageReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+botToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var messageResp SlackMessageResponse
	if err := json.Unmarshal(body, &messageResp); err != nil {
		return err
	}

	if !messageResp.OK {
		return errors.New(fmt.Sprintf("failed to send slack message: %s", messageResp.Error))
	}

	return nil
}

func (s *SlackOAuthService) GetUserFromDB(slackID string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("slack_id = ?", slackID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type SlackMessage struct {
	Type      string `json:"type"`
	User      string `json:"user"`
	BotID     string `json:"bot_id"`
	Subtype   string `json:"subtype"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
}

type SlackMessagesResponse struct {
	OK       bool           `json:"ok"`
	Messages []SlackMessage `json:"messages"`
	Error    string         `json:"error"`
}

func (s *SlackOAuthService) GetChannelMessages(botToken, channelID, oldest string) (*SlackMessagesResponse, error) {
	url := fmt.Sprintf("https://slack.com/api/conversations.history?channel=%s&oldest=%s&limit=10", channelID, oldest)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+botToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var messagesResp SlackMessagesResponse
	if err := json.Unmarshal(body, &messagesResp); err != nil {
		return nil, err
	}

	if !messagesResp.OK {
		return nil, errors.New(fmt.Sprintf("slack api error: %s", messagesResp.Error))
	}

	return &messagesResp, nil
}
