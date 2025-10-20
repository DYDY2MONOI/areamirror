package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

type GoogleDriveService struct {
    client *drive.Service
    config *oauth2.Config
}

func NewGoogleDriveService() (*GoogleDriveService, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("Google OAuth credentials not configured")
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			drive.DriveReadonlyScope,
			drive.DriveMetadataReadonlyScope,
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleDriveService{client: nil, config: cfg}, nil
}

func (g *GoogleDriveService) ListRecentFilesInFolder(userID uint, folderID string, updatedAfter time.Time, pageSize int64) ([]*drive.File, error) {
	if folderID == "" {
		return nil, fmt.Errorf("folderID is required")
	}

	if pageSize <= 0 {
		pageSize = 50
	}

	var oauth2Token models.OAuth2Token
	if err := database.DB.Where("user_id = ? AND service = ?", userID, "google").First(&oauth2Token).Error; err != nil {
		if err := database.DB.Where("user_id = ? AND service = ?", userID, "gmail").First(&oauth2Token).Error; err != nil {
			return nil, fmt.Errorf("no Google OAuth2 token found for user %d: %v", userID, err)
		}
	}

	if oauth2Token.ExpiresAt != nil && time.Now().After(*oauth2Token.ExpiresAt) {
		log.Printf("Google token expired for user %d, attempting refresh", userID)
	}

	token := &oauth2.Token{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		TokenType:    oauth2Token.TokenType,
	}
	if oauth2Token.ExpiresAt != nil {
		token.Expiry = *oauth2Token.ExpiresAt
	}

	ctx := context.Background()
	client := g.config.Client(ctx, token)
	driveService, err := drive.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %v", err)
	}

	q := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	if !updatedAfter.IsZero() {
		q = fmt.Sprintf("%s and modifiedTime > '%s'", q, updatedAfter.UTC().Format(time.RFC3339))
	}

	log.Printf("Drive API query for user %d: %s", userID, q)

	call := driveService.Files.List().Q(q).
		Fields("files(id, name, mimeType, modifiedTime, createdTime, webViewLink, size, owners(displayName, emailAddress))").
		OrderBy("modifiedTime desc").
		PageSize(pageSize)

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("drive list failed: %w", err)
	}

	log.Printf("Drive API returned %d files for folder %s (user %d)", len(resp.Files), folderID, userID)
	for i, file := range resp.Files {
		if i < 3 {
			log.Printf("File %d: ID=%s, Name=%s, Modified=%s", i+1, file.Id, file.Name, file.ModifiedTime)
		}
	}

	return resp.Files, nil
}


