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
	params.Add("scope", "User.Read Files.ReadWrite offline_access")
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

func (o *OneDriveService) RefreshAccessToken(refreshToken string) (*OneDriveTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", o.clientID)
	data.Set("client_secret", o.clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", microsoftTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OneDriveTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	return &tokenResp, nil
}

func (o *OneDriveService) ListFiles(accessToken, folderID string) (*OneDriveListResponse, error) {
	var apiURL string
	if folderID == "" || folderID == "root" {
		apiURL = graphAPIBaseURL + "/me/drive/root/children"
	} else {
		apiURL = fmt.Sprintf("%s/me/drive/items/%s/children", graphAPIBaseURL, folderID)
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list files failed with status %d: %s", resp.StatusCode, string(body))
	}

	var listResp OneDriveListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode list response: %w", err)
	}

	return &listResp, nil
}

func (o *OneDriveService) UploadFile(accessToken, fileName string, content []byte) (*OneDriveUploadResponse, error) {
	apiURL := fmt.Sprintf("%s/me/drive/root:/%s:/content", graphAPIBaseURL, url.PathEscape(fileName))

	req, err := http.NewRequest("PUT", apiURL, bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	var uploadResp OneDriveUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode upload response: %w", err)
	}

	return &uploadResp, nil
}

func (o *OneDriveService) DownloadFile(accessToken, fileID string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/me/drive/items/%s/content", graphAPIBaseURL, fileID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return content, nil
}

func (o *OneDriveService) DeleteFile(accessToken, fileID string) error {
	apiURL := fmt.Sprintf("%s/me/drive/items/%s", graphAPIBaseURL, fileID)

	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (o *OneDriveService) CreateFolder(accessToken, folderName string) (*OneDriveFile, error) {
	apiURL := graphAPIBaseURL + "/me/drive/root/children"

	payload := map[string]interface{}{
		"name":                              folderName,
		"folder":                            map[string]interface{}{},
		"@microsoft.graph.conflictBehavior": "rename",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal folder payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create folder request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create folder failed with status %d: %s", resp.StatusCode, string(body))
	}

	var folder OneDriveFile
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return nil, fmt.Errorf("failed to decode folder response: %w", err)
	}

	return &folder, nil
}

func (o *OneDriveService) GetUserInfo(accessToken string) (map[string]interface{}, error) {
	apiURL := graphAPIBaseURL + "/me"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get user info failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return userInfo, nil
}
