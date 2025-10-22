package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	microsoftAuthURL  = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	microsoftTokenURL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	graphAPIBaseURL   = "https://graph.microsoft.com/v1.0"
)

type OneDriveService struct {
	clientID     string
	clientSecret string
	redirectURI  string
	httpClient   *http.Client
}

type OneDriveTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type OneDriveFile struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Size             int64     `json:"size"`
	WebURL           string    `json:"webUrl"`
	DownloadURL      string    `json:"@microsoft.graph.downloadUrl,omitempty"`
	CreatedDateTime  time.Time `json:"createdDateTime"`
	ModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	Folder           *struct{} `json:"folder,omitempty"`
	File             *struct {
		MimeType string `json:"mimeType"`
	} `json:"file,omitempty"`
}

type OneDriveListResponse struct {
	Value    []OneDriveFile `json:"value"`
	NextLink string         `json:"@odata.nextLink,omitempty"`
}

type OneDriveUploadResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	WebURL string `json:"webUrl"`
}

func NewOneDriveService() (*OneDriveService, error) {
	clientID := os.Getenv("ONEDRIVE_CLIENT_ID")
	clientSecret := os.Getenv("ONEDRIVE_CLIENT_SECRET")
	redirectURI := os.Getenv("ONEDRIVE_REDIRECT_URI")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("OneDrive credentials not configured")
	}

	if redirectURI == "" {
		redirectURI = "http://localhost:8080/onedrive/callback"
	}

	return &OneDriveService{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (o *OneDriveService) GetAuthorizationURL(state string) string {
	params := url.Values{}
	params.Add("client_id", o.clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", o.redirectURI)
	params.Add("scope", "User.Read Files.ReadWrite Files.ReadWrite.All offline_access")
	params.Add("state", state)

	return microsoftAuthURL + "?" + params.Encode()
}

func (o *OneDriveService) ExchangeCodeForToken(code string) (*OneDriveTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", o.clientID)
	data.Set("client_secret", o.clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", o.redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", microsoftTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OneDriveTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}
