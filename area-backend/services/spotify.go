package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	spotifyTokenURL            = "https://accounts.spotify.com/api/token"
	spotifyCurrentlyPlayingURL = "https://api.spotify.com/v1/me/player/currently-playing"
)

type SpotifyService struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
}

var ErrSpotifyPermissionsMissing = errors.New("spotify permissions missing")

type SpotifyAPIError struct {
	Status         int
	Message        string
	Raw            string
	RequiresReauth bool
}

func (e *SpotifyAPIError) Error() string {
	if e == nil {
		return "spotify api error"
	}
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("spotify api error (status %d)", e.Status)
}

type SpotifyNowPlaying struct {
	TrackID       string
	TrackName     string
	Artists       []string
	AlbumName     string
	TrackURL      string
	PreviewURL    string
	DeviceName    string
	IsPlaying     bool
	DurationMS    int
	ProgressMS    int
	StartedAt     time.Time
	CoverImageURL string
}

func NewSpotifyService() (*SpotifyService, error) {
	clientID := strings.TrimSpace(os.Getenv("SPOTIFY_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("SPOTIFY_CLIENT_SECRET"))

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("spotify client credentials are not configured")
	}

	return &SpotifyService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

func (s *SpotifyService) GetCurrentlyPlaying(userID uint) (*SpotifyNowPlaying, error) {
	token, err := s.getSpotifyToken(userID)
	if err != nil {
		return nil, err
	}

	if token.NeedsRefresh() {
		if err := s.refreshSpotifyToken(token); err != nil {
			return nil, fmt.Errorf("failed to refresh spotify token: %w", err)
		}
	}

	return s.fetchCurrentlyPlaying(token, true)
}

func (s *SpotifyService) getSpotifyToken(userID uint) (*models.OAuth2Token, error) {
	var token models.OAuth2Token
	if err := database.DB.Where("user_id = ? AND service = ?", userID, "spotify").First(&token).Error; err != nil {
		return nil, fmt.Errorf("spotify token not found for user %d: %w", userID, err)
	}
	return &token, nil
}

func (s *SpotifyService) refreshSpotifyToken(token *models.OAuth2Token) error {
	refreshToken := strings.TrimSpace(token.RefreshToken)
	if refreshToken == "" {
		return fmt.Errorf("spotify refresh token missing")
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, spotifyTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create spotify refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.clientID, s.clientSecret)))
	req.Header.Set("Authorization", "Basic "+authHeader)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh spotify token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read spotify refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("spotify refresh request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse spotify refresh response: %w", err)
	}

	token.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		token.RefreshToken = tokenResp.RefreshToken
	}
	if tokenResp.TokenType != "" {
		token.TokenType = tokenResp.TokenType
	}
	if tokenResp.Scope != "" {
		token.Scope = tokenResp.Scope
	}

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		token.ExpiresAt = &expiry
	}

	if err := database.DB.Save(token).Error; err != nil {
		return fmt.Errorf("failed to persist refreshed spotify token: %w", err)
	}

	return nil
}

func (s *SpotifyService) fetchCurrentlyPlaying(token *models.OAuth2Token, allowRetry bool) (*SpotifyNowPlaying, error) {
	req, err := http.NewRequest(http.MethodGet, spotifyCurrentlyPlayingURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create spotify currently playing request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query spotify currently playing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read spotify currently playing response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized && allowRetry {
		if err := s.refreshSpotifyToken(token); err != nil {
			return nil, fmt.Errorf("failed to refresh spotify token after unauthorized: %w", err)
		}
		return s.fetchCurrentlyPlaying(token, false)
	}

	if resp.StatusCode >= 300 {
		status := resp.StatusCode
		rawBody := strings.TrimSpace(string(body))
		message := rawBody
		requiresReauth := false

		var apiErr struct {
			Error struct {
				Status  int    `json:"status"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(body, &apiErr); err == nil {
			if apiErr.Error.Status != 0 {
				status = apiErr.Error.Status
			}
			if apiErr.Error.Message != "" {
				message = apiErr.Error.Message
			}
		}

		if status == http.StatusUnauthorized && strings.Contains(strings.ToLower(message), "permission") {
			requiresReauth = true
		}

		return nil, &SpotifyAPIError{
			Status:         status,
			Message:        message,
			Raw:            rawBody,
			RequiresReauth: requiresReauth,
		}
	}

	var payload struct {
		Timestamp            int64  `json:"timestamp"`
		ProgressMs           int    `json:"progress_ms"`
		IsPlaying            bool   `json:"is_playing"`
		CurrentlyPlayingType string `json:"currently_playing_type"`
		Item                 *struct {
			ID           string            `json:"id"`
			Name         string            `json:"name"`
			DurationMs   int               `json:"duration_ms"`
			PreviewURL   string            `json:"preview_url"`
			ExternalURLs map[string]string `json:"external_urls"`
			Artists      []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name   string `json:"name"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
			} `json:"album"`
		} `json:"item"`
		Device *struct {
			Name string `json:"name"`
		} `json:"device"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse spotify currently playing response: %w", err)
	}

	if payload.Item == nil || strings.TrimSpace(payload.Item.ID) == "" {
		return nil, nil
	}

	if payload.CurrentlyPlayingType != "" && payload.CurrentlyPlayingType != "track" {
		return nil, nil
	}

	artists := make([]string, 0, len(payload.Item.Artists))
	for _, artist := range payload.Item.Artists {
		if name := strings.TrimSpace(artist.Name); name != "" {
			artists = append(artists, name)
		}
	}

	var coverURL string
	if len(payload.Item.Album.Images) > 0 {
		coverURL = payload.Item.Album.Images[0].URL
	}

	trackURL := payload.Item.ExternalURLs["spotify"]

	progress := payload.ProgressMs
	if progress < 0 {
		progress = 0
	}

	startedAt := time.Time{}
	if payload.Timestamp > 0 {
		startedAt = time.UnixMilli(payload.Timestamp).Add(-time.Duration(progress) * time.Millisecond)
	}

	deviceName := ""
	if payload.Device != nil {
		deviceName = strings.TrimSpace(payload.Device.Name)
	}

	return &SpotifyNowPlaying{
		TrackID:       payload.Item.ID,
		TrackName:     payload.Item.Name,
		Artists:       artists,
		AlbumName:     payload.Item.Album.Name,
		TrackURL:      trackURL,
		PreviewURL:    payload.Item.PreviewURL,
		DeviceName:    deviceName,
		IsPlaying:     payload.IsPlaying,
		DurationMS:    payload.Item.DurationMs,
		ProgressMS:    progress,
		StartedAt:     startedAt,
		CoverImageURL: coverURL,
	}, nil
}
