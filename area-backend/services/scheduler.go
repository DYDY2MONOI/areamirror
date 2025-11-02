package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type SchedulerService struct {
	emailService    *EmailService
	discordService  *DiscordService
	weatherService  *WeatherService
	sheetsService   *GoogleSheetsService
	driveService    *GoogleDriveService
	telegramService *TelegramService
	openaiService   *OpenAIService
	spotifyService  *SpotifyService
	twitterService  *TwitterService
}

type googleSheetsTriggerConfig struct {
	SpreadsheetID string     `json:"spreadsheetId"`
	Range         string     `json:"range"`
	SheetName     string     `json:"sheetName"`
	HasHeader     bool       `json:"hasHeader"`
	LastValues    [][]string `json:"lastValues"`
	LastChecksum  string     `json:"lastChecksum"`
}

type sheetRowChange struct {
	ChangeType string
	RowNumber  int
	Values     []string
}

type spotifyTriggerConfig struct {
	LastTrackID     string `json:"lastTrackId"`
	LastProgressMs  int    `json:"lastProgressMs"`
	HasProgress     bool   `json:"hasProgress"`
	LastTriggeredAt string `json:"lastTriggeredAt"`
}

type twitterTriggerConfig struct {
	MonitorType      string                       `json:"monitorType"`
	Keyword          string                       `json:"keyword"`
	LastItemID       string                       `json:"lastItemId"`
	IncludeRetweets  bool                         `json:"includeRetweets"`
	TargetUsername   string                       `json:"targetUsername"`
	TargetUserID     string                       `json:"targetUserId"`
	LastTweetStats   map[string]twitterTweetStats `json:"lastTweetStats"`
	KnownFollowerIDs map[string]bool              `json:"knownFollowerIds"`
}

type twitterTweetStats struct {
	LikeCount    int `json:"likeCount"`
	RetweetCount int `json:"retweetCount"`
}

func NewSchedulerService() (*SchedulerService, error) {
	emailService, err := NewEmailService()
	if err != nil {
		log.Printf("Warning: Email service not available: %v", err)
	}

	discordService, err := NewDiscordService()
	if err != nil {
		log.Printf("Warning: Discord service not available: %v", err)
	}

	weatherService, err := NewWeatherService()
	if err != nil {
		log.Printf("Warning: Weather service not available: %v", err)
	}

	sheetsService, err := NewGoogleSheetsService()
	if err != nil {
		log.Printf("Warning: Google Sheets service not available: %v", err)
	}

	telegramService, err := NewTelegramService()
	if err != nil {
		log.Printf("Warning: Telegram service not available: %v", err)
	}

	driveService, err := NewGoogleDriveService()
	if err != nil {
		log.Printf("Warning: Google Drive service not available: %v", err)
	}

	openaiService, err := NewOpenAIService()
	if err != nil {
		log.Printf("Warning: OpenAI service not available: %v", err)
	}

	spotifyService, err := NewSpotifyService()
	if err != nil {
		log.Printf("Warning: Spotify service not available: %v", err)
	}

	twitterService, err := NewTwitterService()
	if err != nil {
		log.Printf("Warning: Twitter service not available: %v", err)
	}

	return &SchedulerService{
		emailService:    emailService,
		discordService:  discordService,
		weatherService:  weatherService,
		sheetsService:   sheetsService,
		driveService:    driveService,
		telegramService: telegramService,
		openaiService:   openaiService,
		spotifyService:  spotifyService,
		twitterService:  twitterService,
	}, nil
}

func (s *SchedulerService) CheckScheduledAreas() error {
	if err := s.checkCalendarTriggers(); err != nil {
		log.Printf("Error checking calendar triggers: %v", err)
	}

	if err := s.checkWeatherTriggers(); err != nil {
		log.Printf("Error checking weather triggers: %v", err)
	}

	if err := s.checkGoogleSheetsTriggers(); err != nil {
		log.Printf("Error checking Google Sheets triggers: %v", err)
	}

	if err := s.checkGoogleDriveTriggers(); err != nil {
		log.Printf("Error checking Google Drive triggers: %v", err)
	}

	if err := s.checkSpotifyTriggers(); err != nil {
		log.Printf("Error checking Spotify triggers: %v", err)
	}

	if err := s.checkTwitterTriggers(); err != nil {
		log.Printf("Error checking Twitter triggers: %v", err)
	}

	if err := s.checkTimerTriggers(); err != nil {
		log.Printf("Error checking timer triggers: %v", err)
	}

	return nil
}

type googleDriveTriggerConfig struct {
	FolderID     string          `json:"folderId"`
	KnownFileIDs map[string]bool `json:"knownFileIds"`
	LastChecked  *time.Time      `json:"lastChecked"`
}

func (s *SchedulerService) checkGoogleDriveTriggers() error {
	if s.driveService == nil {
		log.Printf("Drive service is nil, skipping Drive triggers")
		return nil
	}

	var areas []models.Area
	if err := database.DB.Where("trigger_service = ? AND is_active = ?", "Google Drive", true).Find(&areas).Error; err != nil {
		return fmt.Errorf("failed to fetch google drive areas: %v", err)
	}

	log.Printf("Found %d Google Drive areas to check", len(areas))

	for _, area := range areas {
		log.Printf("Checking Drive area: %s", area.Name)
		var cfg googleDriveTriggerConfig
		if len(area.TriggerConfig) > 0 {
			if err := json.Unmarshal(area.TriggerConfig, &cfg); err != nil {
				log.Printf("Failed to parse Google Drive trigger config for area %s: %v", area.Name, err)
				continue
			}
		}

		if strings.TrimSpace(cfg.FolderID) == "" {
			log.Printf("Google Drive trigger for area %s missing folderId", area.Name)
			continue
		}

		log.Printf("Checking folder %s for area %s", cfg.FolderID, area.Name)

		if cfg.KnownFileIDs == nil {
			cfg.KnownFileIDs = make(map[string]bool)
		}

		files, err := s.driveService.ListRecentFilesInFolder(area.UserID, cfg.FolderID, time.Time{}, 50)
		if err != nil {
			log.Printf("Failed to list drive files for area %s: %v", area.Name, err)
			continue
		}

		log.Printf("Found %d files in folder %s for area %s", len(files), cfg.FolderID, area.Name)

		for _, f := range files {
			if f == nil || f.Id == "" {
				continue
			}
			if cfg.KnownFileIDs[f.Id] {
				continue
			}

			log.Printf("New file detected: %s (%s) in area %s", f.Name, f.Id, area.Name)

			metadata := map[string]interface{}{
				"fileId":       f.Id,
				"fileName":     f.Name,
				"mimeType":     f.MimeType,
				"webViewLink":  f.WebViewLink,
				"createdTime":  f.CreatedTime,
				"modifiedTime": f.ModifiedTime,
				"size":         f.Size,
				"folderId":     cfg.FolderID,
			}

			if err := s.executeArea(area, metadata); err != nil {
				log.Printf("Failed to execute area %s for new Drive file: %v", area.Name, err)
				continue
			}
			cfg.KnownFileIDs[f.Id] = true
		}

		now := time.Now().UTC()
		cfg.LastChecked = &now
		cfgBytes, _ := json.Marshal(cfg)
		if err := database.DB.Model(&area).Update("trigger_config", datatypes.JSON(cfgBytes)).Error; err != nil {
			log.Printf("Failed to persist Drive trigger state for area %s: %v", area.Name, err)
		}
	}

	if err := s.checkTimerTriggers(); err != nil {
		log.Printf("Error checking timer triggers: %v", err)
	}

	return nil
}

func (s *SchedulerService) checkSpotifyTriggers() error {
	if s.spotifyService == nil {
		return nil
	}

	var areas []models.Area
	if err := database.DB.Where("trigger_service = ? AND is_active = ?", "Spotify", true).Find(&areas).Error; err != nil {
		return fmt.Errorf("failed to fetch spotify areas: %v", err)
	}

	for _, area := range areas {
		nowPlaying, err := s.spotifyService.GetCurrentlyPlaying(area.UserID)
		if err != nil {
			if apiErr, ok := err.(*SpotifyAPIError); ok && apiErr.RequiresReauth {
				log.Printf("Spotify permissions missing for user %d (area %s). User must relink Spotify.", area.UserID, area.Name)
			} else {
				log.Printf("Failed to fetch Spotify playback for area %s: %v", area.Name, err)
			}
			continue
		}

		if nowPlaying == nil || !nowPlaying.IsPlaying {
			continue
		}

		var cfg spotifyTriggerConfig
		if len(area.TriggerConfig) > 0 {
			if err := json.Unmarshal(area.TriggerConfig, &cfg); err != nil {
				log.Printf("Failed to parse Spotify trigger config for area %s: %v", area.Name, err)
				cfg = spotifyTriggerConfig{}
			}
		}

		progress := nowPlaying.ProgressMS
		if progress < 0 {
			progress = 0
		}

		shouldTrigger := false
		if strings.TrimSpace(cfg.LastTrackID) == "" {
			shouldTrigger = true
		} else if cfg.LastTrackID != nowPlaying.TrackID {
			shouldTrigger = true
		} else if cfg.HasProgress && progress <= 5000 && cfg.LastProgressMs > 5000 {
			shouldTrigger = true
		}

		cfg.LastTrackID = nowPlaying.TrackID
		cfg.LastProgressMs = progress
		cfg.HasProgress = true

		if shouldTrigger {
			artistNames := strings.Join(nowPlaying.Artists, ", ")
			startedAt := ""
			if !nowPlaying.StartedAt.IsZero() {
				startedAt = nowPlaying.StartedAt.UTC().Format(time.RFC3339)
			}

			metadata := map[string]interface{}{
				"trackId":       nowPlaying.TrackID,
				"trackName":     nowPlaying.TrackName,
				"artistNames":   artistNames,
				"albumName":     nowPlaying.AlbumName,
				"trackUrl":      nowPlaying.TrackURL,
				"previewUrl":    nowPlaying.PreviewURL,
				"deviceName":    nowPlaying.DeviceName,
				"isPlaying":     nowPlaying.IsPlaying,
				"progressMs":    nowPlaying.ProgressMS,
				"durationMs":    nowPlaying.DurationMS,
				"coverImageUrl": nowPlaying.CoverImageURL,
				"eventTitle":    nowPlaying.TrackName,
			}

			if startedAt != "" {
				metadata["startedAt"] = startedAt
				metadata["eventTime"] = startedAt
			}

			if err := s.executeArea(area, metadata); err != nil {
				log.Printf("Failed to execute Spotify area %s: %v", area.Name, err)
			} else {
				cfg.LastTriggeredAt = time.Now().UTC().Format(time.RFC3339)
			}
		}

		cfgBytes, err := json.Marshal(cfg)
		if err != nil {
			log.Printf("Failed to marshal Spotify trigger config for area %s: %v", area.Name, err)
			continue
		}

		if err := database.DB.Model(&area).Update("trigger_config", datatypes.JSON(cfgBytes)).Error; err != nil {
			log.Printf("Failed to persist Spotify trigger config for area %s: %v", area.Name, err)
		}
	}

	return nil
}

func (s *SchedulerService) checkTwitterTriggers() error {
	if s.twitterService == nil {
		return nil
	}

	var areas []models.Area
	if err := database.DB.Where("trigger_service = ? AND is_active = ?", "Twitter", true).Find(&areas).Error; err != nil {
		return fmt.Errorf("failed to fetch twitter areas: %v", err)
	}

	for _, area := range areas {
		var cfg twitterTriggerConfig
		if len(area.TriggerConfig) > 0 {
			if err := json.Unmarshal(area.TriggerConfig, &cfg); err != nil {
				log.Printf("Failed to parse Twitter trigger config for area %s: %v", area.Name, err)
				cfg = twitterTriggerConfig{}
			}
		}

		monitorType := strings.ToLower(strings.TrimSpace(cfg.MonitorType))
		if monitorType == "" {
			monitorType = "mentions"
		}
		cfg.MonitorType = monitorType

		skipArea := false
		handleAPIError := func(err error, context string) {
			skipArea = true
			if apiErr, ok := err.(*TwitterAPIError); ok {
				switch apiErr.StatusCode {
				case http.StatusTooManyRequests:
					log.Printf("Twitter rate limit for area %s (%s); will retry on next cycle", area.Name, context)
				case http.StatusForbidden:
					errMsg := apiErr.Message
					if strings.TrimSpace(errMsg) == "" {
						errMsg = "Twitter API returned 403 Forbidden"
					}
					s.recordAreaFailure(&area, fmt.Errorf("twitter access forbidden for %s: %s", context, errMsg))
				default:
					log.Printf("Twitter API error while %s for area %s: %v", context, area.Name, err)
				}
				return
			}
			log.Printf("Failed to %s for area %s: %v", context, area.Name, err)
		}

		var user models.User
		if err := database.DB.First(&user, area.UserID).Error; err != nil {
			log.Printf("Failed to load user %d for twitter area %s: %v", area.UserID, area.Name, err)
			continue
		}

		twitterUserID := strings.TrimSpace(cfg.TargetUserID)
		if twitterUserID == "" && user.TwitterID != nil {
			twitterUserID = strings.TrimSpace(*user.TwitterID)
		}
		if twitterUserID == "" {
			log.Printf("Skipping twitter area %s: user %d has no linked twitter account", area.Name, area.UserID)
			continue
		}

		cfgDirty := false
		if cfg.TargetUserID == "" {
			cfg.TargetUserID = twitterUserID
			cfgDirty = true
		}
		if cfg.TargetUsername == "" && user.TwitterUsername != nil {
			cfg.TargetUsername = *user.TwitterUsername
			cfgDirty = true
		}

		switch monitorType {
		case "followers":
			if cfg.KnownFollowerIDs == nil {
				cfg.KnownFollowerIDs = make(map[string]bool)
				cfgDirty = true
			}

			followers, err := s.twitterService.FetchFollowers(area.UserID, twitterUserID, 100)
			if err != nil {
				handleAPIError(err, "fetch twitter followers")
				break
			}

			for _, follower := range followers {
				if cfg.KnownFollowerIDs[follower.ID] {
					continue
				}

				metadata := map[string]interface{}{
					"followerId":         follower.ID,
					"followerUsername":   follower.Username,
					"followerName":       follower.Name,
					"followerBio":        follower.Description,
					"followerCreatedAt":  follower.CreatedAt.Format(time.RFC3339),
					"twitterMonitorType": monitorType,
				}
				if user.TwitterUsername != nil {
					metadata["accountUsername"] = *user.TwitterUsername
				}

				if err := s.executeArea(area, metadata); err != nil {
					log.Printf("Failed to execute twitter follower area %s for follower %s: %v", area.Name, follower.ID, err)
					continue
				}

				cfg.KnownFollowerIDs[follower.ID] = true
				cfgDirty = true
			}

			if len(cfg.KnownFollowerIDs) > 2000 && len(followers) > 0 {
				trimmed := make(map[string]bool, len(followers))
				for _, follower := range followers {
					trimmed[follower.ID] = true
				}
				cfg.KnownFollowerIDs = trimmed
				cfgDirty = true
			}
		case "likes", "retweets":
			if cfg.LastTweetStats == nil {
				cfg.LastTweetStats = make(map[string]twitterTweetStats)
				cfgDirty = true
			}

			tweets, err := s.twitterService.FetchUserTweetsWithMetrics(area.UserID, twitterUserID, 50)
			if err != nil {
				handleAPIError(err, "fetch twitter tweet metrics")
				break
			}

			currentIDs := make(map[string]bool, len(tweets))
			username := ""
			if user.TwitterUsername != nil {
				username = *user.TwitterUsername
			}

			fullName := strings.TrimSpace(strings.TrimSpace(user.FirstName + " " + user.LastName))
			for _, tweet := range tweets {
				currentIDs[tweet.ID] = true
				prev := cfg.LastTweetStats[tweet.ID]

				if monitorType == "likes" {
					newLikes := tweet.LikeCount - prev.LikeCount
					if newLikes > 0 {
						metadata := map[string]interface{}{
							"tweetId":            tweet.ID,
							"tweetText":          tweet.Text,
							"tweetCreatedAt":     tweet.CreatedAt.Format(time.RFC3339),
							"tweetLikeCount":     tweet.LikeCount,
							"tweetNewLikes":      newLikes,
							"twitterMonitorType": monitorType,
							"eventTitle":         tweet.Text,
						}
						if username != "" {
							metadata["accountUsername"] = username
							metadata["tweetAuthorUsername"] = username
							metadata["tweetUrl"] = buildTweetURL(username, tweet.ID)
						}
						if user.TwitterID != nil {
							metadata["tweetAuthorId"] = *user.TwitterID
						}
						if fullName != "" {
							metadata["tweetAuthorName"] = fullName
						}
						if err := s.executeArea(area, metadata); err != nil {
							log.Printf("Failed to execute twitter like area %s for tweet %s: %v", area.Name, tweet.ID, err)
						}
					}
				} else if monitorType == "retweets" {
					newRetweets := tweet.RetweetCount - prev.RetweetCount
					if newRetweets > 0 {
						metadata := map[string]interface{}{
							"tweetId":            tweet.ID,
							"tweetText":          tweet.Text,
							"tweetCreatedAt":     tweet.CreatedAt.Format(time.RFC3339),
							"tweetRetweetCount":  tweet.RetweetCount,
							"tweetNewRetweets":   newRetweets,
							"twitterMonitorType": monitorType,
							"eventTitle":         tweet.Text,
						}
						if username != "" {
							metadata["accountUsername"] = username
							metadata["tweetAuthorUsername"] = username
							metadata["tweetUrl"] = buildTweetURL(username, tweet.ID)
						}
						if user.TwitterID != nil {
							metadata["tweetAuthorId"] = *user.TwitterID
						}
						if fullName != "" {
							metadata["tweetAuthorName"] = fullName
						}
						if err := s.executeArea(area, metadata); err != nil {
							log.Printf("Failed to execute twitter retweet area %s for tweet %s: %v", area.Name, tweet.ID, err)
						}
					}
				}

				cfg.LastTweetStats[tweet.ID] = twitterTweetStats{
					LikeCount:    tweet.LikeCount,
					RetweetCount: tweet.RetweetCount,
				}
				cfgDirty = true
			}

			// avoid unbounded growth: keep only tweets seen recently
			if len(cfg.LastTweetStats) > 0 {
				for id := range cfg.LastTweetStats {
					if !currentIDs[id] && len(cfg.LastTweetStats) > 200 {
						delete(cfg.LastTweetStats, id)
						cfgDirty = true
					}
				}
			}
		default: // mentions
			tweets, newestID, err := s.twitterService.FetchMentions(area.UserID, twitterUserID, cfg.LastItemID)
			if err != nil {
				handleAPIError(err, "fetch twitter mentions")
				break
			}

			if len(tweets) == 0 {
				if newestID != "" && newestID != cfg.LastItemID {
					cfg.LastItemID = newestID
					cfgDirty = true
				}
				break
			}

			keyword := strings.ToLower(strings.TrimSpace(cfg.Keyword))
			includeRetweets := cfg.IncludeRetweets
			lastProcessedID := cfg.LastItemID

			for _, tweet := range tweets {
				lastProcessedID = tweet.ID

				if !includeRetweets && strings.HasPrefix(strings.TrimSpace(tweet.Text), "RT ") {
					continue
				}

				if keyword != "" && !strings.Contains(strings.ToLower(tweet.Text), keyword) {
					continue
				}

				tweetURL := tweet.URL
				if tweetURL == "" {
					username := tweet.AuthorUsername
					if username == "" && user.TwitterUsername != nil {
						username = *user.TwitterUsername
					}
					tweetURL = buildTweetURL(username, tweet.ID)
				}

				metadata := map[string]interface{}{
					"tweetId":             tweet.ID,
					"tweetText":           tweet.Text,
					"tweetAuthorId":       tweet.AuthorID,
					"tweetAuthorUsername": tweet.AuthorUsername,
					"tweetAuthorName":     tweet.AuthorName,
					"tweetCreatedAt":      tweet.CreatedAt.Format(time.RFC3339),
					"tweetUrl":            tweetURL,
					"twitterMonitorType":  monitorType,
					"eventTitle":          tweet.Text,
				}

				if tweet.InReplyToUserID != "" {
					metadata["tweetInReplyToUserId"] = tweet.InReplyToUserID
				}
				if tweet.ConversationID != "" {
					metadata["tweetConversationId"] = tweet.ConversationID
				}
				if user.TwitterUsername != nil {
					metadata["accountUsername"] = *user.TwitterUsername
				}

				if err := s.executeArea(area, metadata); err != nil {
					log.Printf("Failed to execute twitter area %s for tweet %s: %v", area.Name, tweet.ID, err)
					continue
				}
			}

			if newestID != "" {
				lastProcessedID = newestID
			}

			if lastProcessedID != cfg.LastItemID {
				cfg.LastItemID = lastProcessedID
				cfgDirty = true
			}
		}

		if skipArea {
			if cfgDirty {
				if err := s.persistTwitterConfig(area, cfg); err != nil {
					log.Printf("Failed to persist Twitter config for area %s: %v", area.Name, err)
				}
			}
			continue
		}

		if cfgDirty {
			if err := s.persistTwitterConfig(area, cfg); err != nil {
				log.Printf("Failed to persist Twitter config for area %s: %v", area.Name, err)
			}
		}
	}

	return nil
}

func (s *SchedulerService) checkCalendarTriggers() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service IN ? AND is_active = ?", []string{"Google Calendar", "Date Timer"}, true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch calendar areas: %v", err)
	}

	now := time.Now()

	for _, area := range areas {
		var triggerConfig map[string]interface{}
		if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
			log.Printf("Failed to parse trigger config for area %s: %v", area.Name, err)
			continue
		}

		if s.shouldTriggerArea(area, triggerConfig, now) {
			if err := s.executeArea(area, nil); err != nil {
				log.Printf("Failed to execute area %s: %v", area.Name, err)
			}
		}
	}

	return nil
}

func (s *SchedulerService) checkWeatherTriggers() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Weather", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch weather areas: %v", err)
	}

	for _, area := range areas {
		var triggerConfig map[string]interface{}
		if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
			log.Printf("Failed to parse trigger config for area %s: %v", area.Name, err)
			continue
		}

		if s.shouldTriggerWeatherArea(area, triggerConfig) {
			if err := s.executeArea(area, nil); err != nil {
				log.Printf("Failed to execute area %s: %v", area.Name, err)
			}
		}
	}

	return nil
}

func (s *SchedulerService) checkTimerTriggers() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Timer", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch timer areas: %v", err)
	}

	log.Printf("Checking Timer triggers: found %d active Timer areas", len(areas))

	now := time.Now()

	for _, area := range areas {
		var triggerConfig map[string]interface{}
		if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
			log.Printf("Failed to parse trigger config for area %s: %v", area.Name, err)
			continue
		}

		if s.shouldTriggerTimerArea(area, triggerConfig, now) {
			metadata := map[string]interface{}{
				"triggerTime": now.Format("2006-01-02 15:04:05"),
				"timerName":   area.Name,
			}
			if intervalStr, ok := triggerConfig["interval"].(string); ok {
				metadata["interval"] = intervalStr
			}

			if err := s.executeArea(area, metadata); err != nil {
				log.Printf("Failed to execute timer area %s: %v", area.Name, err)
			}
		}
	}

	return nil
}

func (s *SchedulerService) checkGoogleSheetsTriggers() error {
	if s.sheetsService == nil {
		return nil
	}

	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Google Sheets", true).Find(&areas).Error
	if err != nil {
		return fmt.Errorf("failed to fetch google sheets areas: %v", err)
	}

	for _, area := range areas {
		var cfg googleSheetsTriggerConfig
		if len(area.TriggerConfig) > 0 {
			if err := json.Unmarshal(area.TriggerConfig, &cfg); err != nil {
				log.Printf("Failed to parse Google Sheets trigger config for area %s: %v", area.Name, err)
				continue
			}
		}

		if cfg.SpreadsheetID == "" || cfg.Range == "" {
			log.Printf("Google Sheets trigger for area %s missing spreadsheetId or range", area.Name)
			continue
		}

		rows, err := s.sheetsService.FetchValues(cfg.SpreadsheetID, cfg.Range)
		if err != nil {
			log.Printf("Failed to fetch sheet values for area %s: %v", area.Name, err)
			continue
		}

		headers, dataRows := splitSheetHeader(rows, cfg.HasHeader)

		if cfg.LastChecksum == "" && len(cfg.LastValues) == 0 {
			cfg.LastValues = dataRows
			cfg.LastChecksum = hashSheetRows(dataRows)
			if err := s.persistGoogleSheetsConfig(area, cfg); err != nil {
				log.Printf("Failed to persist initial sheet snapshot for area %s: %v", area.Name, err)
			}
			continue
		}

		change, ok := detectSheetChange(cfg.LastValues, dataRows)
		if !ok {
			continue
		}

		rowNumber := change.RowNumber
		if cfg.HasHeader {
			rowNumber++
		}

		rowData := buildRowMap(headers, change.Values)
		metadata := map[string]interface{}{
			"spreadsheetId":  cfg.SpreadsheetID,
			"spreadsheetUrl": fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", cfg.SpreadsheetID),
			"sheetName":      cfg.SheetName,
			"changeType":     change.ChangeType,
			"rowNumber":      rowNumber,
			"rowValues":      change.Values,
			"rowData":        rowData,
		}

		if err := s.executeArea(area, metadata); err != nil {
			log.Printf("Failed to execute area %s: %v", area.Name, err)
			continue
		}

		cfg.LastValues = dataRows
		cfg.LastChecksum = hashSheetRows(dataRows)
		if err := s.persistGoogleSheetsConfig(area, cfg); err != nil {
			log.Printf("Failed to update sheet snapshot for area %s: %v", area.Name, err)
		}
	}

	return nil
}

func (s *SchedulerService) shouldTriggerArea(area models.Area, triggerConfig map[string]interface{}, now time.Time) bool {
	eventTimeStr, ok := triggerConfig["eventTime"].(string)
	if !ok {
		return false
	}

	eventTime, err := time.Parse(time.RFC3339, eventTimeStr)
	if err != nil {
		log.Printf("Failed to parse event time for area %s: %v", area.Name, err)
		return false
	}

	timeDiff := eventTime.Sub(now)
	return timeDiff >= 0 && timeDiff <= 30*time.Second
}

func (s *SchedulerService) shouldTriggerWeatherArea(area models.Area, triggerConfig map[string]interface{}) bool {
	if s.weatherService == nil {
		log.Printf("Weather service not available for area %s", area.Name)
		return false
	}

	city, ok := triggerConfig["city"].(string)
	if !ok || city == "" {
		log.Printf("City not specified for weather area %s", area.Name)
		return false
	}

	temperature, ok := triggerConfig["temperature"].(float64)
	if !ok {
		temperature = 0
	}

	condition, ok := triggerConfig["condition"].(string)
	if !ok {
		condition = ""
	}

	operator, ok := triggerConfig["operator"].(string)
	if !ok {
		operator = "greater_than"
	}

	weatherConfig := WeatherTriggerConfig{
		City:        city,
		Temperature: temperature,
		Condition:   condition,
		Operator:    operator,
	}

	result, err := s.weatherService.CheckWeatherTrigger(weatherConfig)
	if err != nil {
		log.Printf("Failed to check weather trigger for area %s: %v", area.Name, err)
		return false
	}

	if result.Triggered {
		log.Printf("Weather trigger activated for area %s: %s", area.Name, result.Message)
		return true
	}

	return false
}

func (s *SchedulerService) shouldTriggerTimerArea(area models.Area, triggerConfig map[string]interface{}, now time.Time) bool {
	intervalStr, ok := triggerConfig["interval"].(string)
	if !ok || intervalStr == "" {
		log.Printf("Timer area %s missing interval config", area.Name)
		return false
	}

	interval, err := parseDuration(intervalStr)
	if err != nil {
		log.Printf("Timer area %s has invalid interval format: %v", area.Name, err)
		return false
	}

	if area.LastRunAt != nil {
		timeSinceLastRun := now.Sub(*area.LastRunAt)
		if timeSinceLastRun < interval {
			log.Printf("Timer area %s not yet ready (last run %v ago, need %v)",
				area.Name, timeSinceLastRun, interval)
			return false
		}
	}

	log.Printf("Timer trigger activated for area %s (interval: %s)", area.Name, intervalStr)
	return true
}

func (s *SchedulerService) executeArea(area models.Area, metadata map[string]interface{}) error {
	log.Printf("Executing area: %s", area.Name)
	log.Printf("IntermediateActionService: '%s'", area.IntermediateActionService)

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	if area.IntermediateActionService != "" && area.IntermediateActionService == "OpenAI" {
		log.Printf("Executing OpenAI intermediate action for area: %s", area.Name)
		if err := s.executeIntermediateAction(&area, metadata); err != nil {
			log.Printf("WARNING: Failed to execute intermediate action: %v", err)
			log.Printf("Continuing with action execution even though OpenAI failed")
			metadata["openaiGeneratedText"] = "[Erreur OpenAI: " + err.Error() + "]"
		}
		log.Printf("OpenAI intermediate action completed for area: %s", area.Name)
	} else {
		log.Printf("Skipping intermediate action (empty or not OpenAI): '%s'", area.IntermediateActionService)
	}

	var actionConfig map[string]interface{}
	if err := json.Unmarshal(area.ActionConfig, &actionConfig); err != nil {
		return fmt.Errorf("failed to parse action config: %v", err)
	}

	switch area.ActionService {
	case "Gmail":
		return s.executeGmailAction(&area, actionConfig, metadata)
	case "Discord":
		return s.executeDiscordAction(&area, actionConfig, metadata)
	case "Telegram":
		return s.executeTelegramAction(&area, actionConfig, metadata)
	case "Spotify":
		return s.executeSpotifyAction(&area, actionConfig, metadata)
	case "Twitter":
		return s.executeTwitterAction(&area, actionConfig, metadata)
	default:
		log.Printf("Unsupported action service: %s", area.ActionService)
		return nil
	}
}

func (s *SchedulerService) ExecuteAreaPublic(area models.Area, metadata map[string]interface{}) error {
	return s.executeArea(area, metadata)
}

func (s *SchedulerService) executeIntermediateAction(area *models.Area, metadata map[string]interface{}) error {
	switch area.IntermediateActionService {
	case "OpenAI":
		return s.executeOpenAIAction(area, metadata)
	default:
		return fmt.Errorf("unsupported intermediate action service: %s", area.IntermediateActionService)
	}
}

func (s *SchedulerService) executeOpenAIAction(area *models.Area, metadata map[string]interface{}) error {
	if s.openaiService == nil {
		return fmt.Errorf("OpenAI service not available")
	}

	var intermediateConfig map[string]interface{}
	if err := json.Unmarshal(area.IntermediateActionConfig, &intermediateConfig); err != nil {
		return fmt.Errorf("failed to parse intermediate action config: %v", err)
	}

	prompt := getString(intermediateConfig["prompt"])
	if prompt == "" {
		return fmt.Errorf("prompt not found in intermediate action config")
	}

	systemPrompt := getString(intermediateConfig["systemPrompt"])
	temperature := 0.7
	if temp, ok := intermediateConfig["temperature"].(float64); ok {
		temperature = temp
	}
	maxTokens := 500
	if tokens, ok := intermediateConfig["maxTokens"].(float64); ok {
		maxTokens = int(tokens)
	}

	templateVars := buildTemplateVars(area, metadata)
	prompt = applyTemplateVariables(prompt, templateVars)
	if systemPrompt != "" {
		systemPrompt = applyTemplateVariables(systemPrompt, templateVars)
	}

	generatedText, err := s.openaiService.GenerateText(prompt, systemPrompt, temperature, maxTokens)
	if err != nil {
		return fmt.Errorf("failed to generate text with OpenAI: %w", err)
	}

	metadata["openaiGeneratedText"] = generatedText
	log.Printf("OpenAI generated text (length: %d): %s", len(generatedText), generatedText)
	if len(generatedText) == 0 {
		log.Printf("WARNING: OpenAI generated empty text!")
	}
	return nil
}

func (s *SchedulerService) executeGmailAction(area *models.Area, actionConfig map[string]interface{}, metadata map[string]interface{}) error {
	if s.emailService == nil {
		return fmt.Errorf("Email service not available")
	}

	toEmail := strings.TrimSpace(getString(actionConfig["toEmail"]))
	if toEmail == "" {
		return fmt.Errorf("toEmail not found in action config")
	}

	subject := getString(actionConfig["subject"])
	if strings.TrimSpace(subject) == "" {
		subject = "AREA Notification"
	}

	body := getString(actionConfig["body"])
	if strings.TrimSpace(body) == "" {
		body = "This is an automated message from your AREA."
	}

	templateVars := buildTemplateVars(area, metadata)
	subject = applyTemplateVariables(subject, templateVars)
	body = applyTemplateVariables(body, templateVars)

	emailReq := EmailRequest{
		To:      toEmail,
		Subject: subject,
		Body:    body,
	}

	if err := s.emailService.SendEmail(emailReq); err != nil {
		s.recordAreaFailure(area, fmt.Errorf("failed to send email: %w", err))
		return fmt.Errorf("failed to send email: %v", err)
	}

	s.recordAreaSuccess(area)
	log.Printf("Email sent successfully for AREA: %s", area.Name)
	log.Printf("Successfully executed area: %s", area.Name)
	return nil
}

func (s *SchedulerService) executeDiscordAction(area *models.Area, actionConfig map[string]interface{}, metadata map[string]interface{}) error {
	if s.discordService == nil {
		return fmt.Errorf("Discord service not available")
	}

	webhookURL, _ := actionConfig["webhookUrl"].(string)
	if webhookURL == "" {
		if alt, ok := actionConfig["webhookURL"].(string); ok {
			webhookURL = alt
		}
	}
	if webhookURL == "" {
		return fmt.Errorf("webhookUrl not found in action config")
	}

	message, _ := actionConfig["message"].(string)
	if message == "" {
		message = fmt.Sprintf("Notification from area %s", area.Name)
	}

	templateVars := buildTemplateVars(area, metadata)
	message = applyTemplateVariables(message, templateVars)

	if err := s.discordService.SendWebhookMessage(webhookURL, message); err != nil {
		s.recordAreaFailure(area, fmt.Errorf("failed to send discord message: %w", err))
		return fmt.Errorf("failed to send discord message: %v", err)
	}

	s.recordAreaSuccess(area)
	log.Printf("Discord message sent successfully for AREA: %s", area.Name)

	if metadata != nil {
		s.persistDiscordLog(area, message, metadata)
	}

	log.Printf("Successfully executed area: %s", area.Name)
	return nil
}

func (s *SchedulerService) executeTelegramAction(area *models.Area, actionConfig map[string]interface{}, metadata map[string]interface{}) error {
	if s.telegramService == nil {
		return fmt.Errorf("Telegram service not available")
	}

	chatID := strings.TrimSpace(getString(actionConfig["chatId"]))
	if chatID == "" {
		chatID = strings.TrimSpace(getString(actionConfig["chatID"]))
	}
	if chatID == "" {
		chatID = strings.TrimSpace(getString(actionConfig["chat_id"]))
	}
	if chatID == "" {
		return fmt.Errorf("chatId not found in action config")
	}

	message := getString(actionConfig["message"])
	if strings.TrimSpace(message) == "" {
		message = fmt.Sprintf("Notification from area %s", area.Name)
	}

	templateVars := buildTemplateVars(area, metadata)
	log.Printf("Template vars for Telegram: openaiGeneratedText='%s'", templateVars["openaiGeneratedText"])
	log.Printf("Message before template replacement: '%s'", message)
	message = applyTemplateVariables(message, templateVars)
	log.Printf("Message after template replacement: '%s'", message)

	if err := s.telegramService.SendMessage(chatID, message); err != nil {
		s.recordAreaFailure(area, fmt.Errorf("failed to send telegram message: %w", err))
		return fmt.Errorf("failed to send telegram message: %v", err)
	}

	s.recordAreaSuccess(area)
	log.Printf("Telegram message sent successfully for AREA: %s", area.Name)
	log.Printf("Successfully executed area: %s", area.Name)
	return nil
}

func (s *SchedulerService) executeTwitterAction(area *models.Area, actionConfig map[string]interface{}, metadata map[string]interface{}) error {
	if s.twitterService == nil {
		return fmt.Errorf("Twitter service not available")
	}

	templateVars := buildTemplateVars(area, metadata)
	modeValue := strings.ToLower(strings.TrimSpace(getString(actionConfig["actionMode"])))
	if modeValue == "" {
		modeFromArea := strings.ToLower(strings.TrimSpace(area.ActionType))
		switch {
		case strings.Contains(modeFromArea, "retweet"):
			modeValue = "retweet"
		case strings.Contains(modeFromArea, "tweet"):
			modeValue = "tweet"
		default:
			modeValue = "tweet"
		}
	}

	if metadata != nil {
		metadata["twitterActionMode"] = modeValue
	}

	var user models.User
	userErr := database.DB.First(&user, area.UserID).Error

	switch modeValue {
	case "retweet":
		if userErr != nil {
			errMsg := fmt.Errorf("failed to load user for twitter retweet: %w", userErr)
			s.recordAreaFailure(area, errMsg)
			return errMsg
		}
		if user.TwitterID == nil || strings.TrimSpace(*user.TwitterID) == "" {
			errMsg := fmt.Errorf("twitter account not linked for retweet action")
			s.recordAreaFailure(area, errMsg)
			return errMsg
		}

		tweetID := strings.TrimSpace(getString(actionConfig["tweetId"]))
		if tweetID != "" {
			tweetID = applyTemplateVariables(tweetID, templateVars)
		}
		if tweetID == "" && metadata != nil {
			tweetID = strings.TrimSpace(getString(metadata["tweetId"]))
		}
		if tweetID == "" {
			errMsg := fmt.Errorf("tweetId not found in action config or metadata")
			s.recordAreaFailure(area, errMsg)
			return errMsg
		}

		if err := s.twitterService.Retweet(area.UserID, *user.TwitterID, tweetID); err != nil {
			errMsg := fmt.Errorf("failed to retweet tweet: %w", err)
			s.recordAreaFailure(area, errMsg)
			return errMsg
		}

		if metadata != nil {
			metadata["tweetId"] = tweetID
			if user.TwitterUsername != nil {
				metadata["retweetPerformedBy"] = *user.TwitterUsername
				if strings.TrimSpace(getString(metadata["tweetUrl"])) == "" {
					metadata["tweetUrl"] = buildTweetURL(*user.TwitterUsername, tweetID)
				}
			}
		}

		s.recordAreaSuccess(area)
		log.Printf("Retweeted tweet %s for AREA: %s", tweetID, area.Name)
		return nil
	default:
		tweetText := strings.TrimSpace(getString(actionConfig["tweetText"]))
		if tweetText == "" {
			return fmt.Errorf("tweetText not found in action config")
		}

		tweetText = applyTemplateVariables(tweetText, templateVars)

		replyID := strings.TrimSpace(getString(actionConfig["replyToTweetId"]))
		if replyID != "" {
			replyID = applyTemplateVariables(replyID, templateVars)
		}

		tweet, err := s.twitterService.PostTweet(area.UserID, tweetText, replyID)
		if err != nil {
			s.recordAreaFailure(area, fmt.Errorf("failed to post tweet: %w", err))
			return fmt.Errorf("failed to post tweet: %v", err)
		}

		if userErr == nil {
			if user.TwitterUsername != nil && tweet.URL == "" {
				tweet.URL = buildTweetURL(*user.TwitterUsername, tweet.ID)
			}
			if metadata != nil {
				if user.TwitterUsername != nil {
					metadata["tweetAuthorUsername"] = *user.TwitterUsername
					metadata["accountUsername"] = *user.TwitterUsername
				}
				if user.TwitterID != nil {
					metadata["tweetAuthorId"] = *user.TwitterID
				}
				fullName := strings.TrimSpace(strings.TrimSpace(user.FirstName + " " + user.LastName))
				if fullName != "" {
					metadata["tweetAuthorName"] = fullName
				}
			}
		}

		if metadata != nil {
			metadata["tweetId"] = tweet.ID
			metadata["tweetText"] = tweet.Text
			metadata["tweetCreatedAt"] = tweet.CreatedAt.Format(time.RFC3339)
			if tweet.URL != "" {
				metadata["tweetUrl"] = tweet.URL
			}
			if replyID != "" {
				metadata["replyToTweetId"] = replyID
			}
		}

		s.recordAreaSuccess(area)
		log.Printf("Tweet posted successfully for AREA: %s (tweet ID: %s)", area.Name, tweet.ID)
		return nil
	}
}

func (s *SchedulerService) executeSpotifyAction(area *models.Area, actionConfig map[string]interface{}, metadata map[string]interface{}) error {
	if s.spotifyService == nil {
		return fmt.Errorf("Spotify service not available")
	}
	if s.sheetsService == nil {
		return fmt.Errorf("Google Sheets service not available")
	}

	playlistID := strings.TrimSpace(getString(actionConfig["playlistId"]))
	if playlistID == "" {
		playlistID = strings.TrimSpace(getString(actionConfig["playlistID"]))
	}
	if playlistID == "" {
		return fmt.Errorf("playlistId not found in action config")
	}

	spreadsheetID := strings.TrimSpace(getString(actionConfig["spreadsheetId"]))
	if spreadsheetID == "" {
		spreadsheetID = strings.TrimSpace(getString(actionConfig["sheetId"]))
	}
	if spreadsheetID == "" && metadata != nil {
		spreadsheetID = strings.TrimSpace(getString(metadata["spreadsheetId"]))
	}
	if spreadsheetID == "" {
		return fmt.Errorf("spreadsheetId not found in action config or metadata")
	}

	readRange := strings.TrimSpace(getString(actionConfig["range"]))
	if readRange == "" {
		readRange = strings.TrimSpace(getString(actionConfig["sheetRange"]))
	}
	if readRange == "" && metadata != nil {
		if name := strings.TrimSpace(getString(metadata["sheetName"])); name != "" {
			readRange = name
		}
	}
	if readRange == "" {
		return fmt.Errorf("range not found in action config")
	}

	hasHeader := true
	if val, ok := actionConfig["hasHeader"].(bool); ok {
		hasHeader = val
	} else if val, ok := actionConfig["has_header"].(bool); ok {
		hasHeader = val
	}

	columnIndex := -1
	if idx, ok := extractInt(actionConfig["urlColumnIndex"]); ok && idx >= 0 {
		columnIndex = idx
	} else if idx, ok := extractInt(actionConfig["linkColumnIndex"]); ok && idx >= 0 {
		columnIndex = idx
	}

	columnIdentifier := strings.TrimSpace(getString(actionConfig["urlColumn"]))
	if columnIdentifier == "" {
		columnIdentifier = strings.TrimSpace(getString(actionConfig["spotifyColumn"]))
	}
	if columnIdentifier == "" {
		columnIdentifier = "SpotifyLink"
	}

	rows, err := s.sheetsService.FetchValues(spreadsheetID, readRange)
	if err != nil {
		wrapped := fmt.Errorf("failed to fetch sheet values: %w", err)
		s.recordAreaFailure(area, wrapped)
		return wrapped
	}

	if hasHeader && len(rows) == 0 {
		s.recordAreaFailure(area, fmt.Errorf("sheet range %s appears empty", readRange))
		return fmt.Errorf("sheet range %s appears empty", readRange)
	}

	headers := []string{}
	startRow := 0
	if hasHeader && len(rows) > 0 {
		headers = rows[0]
		startRow = 1
	}

	if columnIndex < 0 {
		columnIndex = resolveColumnIndex(headers, columnIdentifier)
	}
	if columnIndex < 0 {
		if idx, err := parseColumnIdentifier(columnIdentifier); err == nil {
			columnIndex = idx
		}
	}
	if columnIndex < 0 && len(headers) > 0 {
		// try fallback header that matches common variants
		if idx := resolveColumnIndex(headers, "spotify url"); idx >= 0 {
			columnIndex = idx
		} else if idx := resolveColumnIndex(headers, "spotify"); idx >= 0 {
			columnIndex = idx
		}
	}
	if columnIndex < 0 {
		wrapped := fmt.Errorf("could not locate Spotify link column '%s'", columnIdentifier)
		s.recordAreaFailure(area, wrapped)
		return wrapped
	}

	trackURIs := make([]string, 0, len(rows))
	for _, row := range rows[startRow:] {
		if columnIndex >= len(row) {
			continue
		}
		rawValue := row[columnIndex]
		uri, err := normalizeSpotifyTrackURI(rawValue)
		if err != nil {
			log.Printf("Skipping invalid Spotify link '%s' for area %s: %v", rawValue, area.Name, err)
			continue
		}
		if uri != "" {
			trackURIs = append(trackURIs, uri)
		}
	}

	if err := s.spotifyService.ReplacePlaylistTracks(area.UserID, playlistID, trackURIs); err != nil {
		wrapped := fmt.Errorf("failed to update Spotify playlist: %w", err)
		s.recordAreaFailure(area, wrapped)
		return wrapped
	}

	s.recordAreaSuccess(area)
	log.Printf("Spotify playlist %s updated with %d tracks for area %s", playlistID, len(trackURIs), area.Name)
	return nil
}

func (s *SchedulerService) recordAreaSuccess(area *models.Area) {
	now := time.Now()
	area.LastRunAt = &time.Time{}
	*area.LastRunAt = now
	area.RunCount++
	area.LastRunStatus = "success"
	area.LastError = ""

	if err := database.DB.Save(area).Error; err != nil {
		log.Printf("Failed to update area status for %s: %v", area.Name, err)
	}
}

func (s *SchedulerService) recordAreaFailure(area *models.Area, runErr error) {
	now := time.Now()
	area.LastRunAt = &time.Time{}
	*area.LastRunAt = now
	area.LastRunStatus = "failed"
	if runErr != nil {
		area.LastError = runErr.Error()
	}

	if err := database.DB.Save(area).Error; err != nil {
		log.Printf("Failed to persist failed status for area %s: %v", area.Name, err)
	}
}

func (s *SchedulerService) persistDiscordLog(area *models.Area, message string, metadata map[string]interface{}) {
	logEntry := models.DiscordMessageLog{
		AreaID:  area.ID,
		Message: message,
	}

	if filePath, ok := metadata["spreadsheetUrl"].(string); ok && filePath != "" {
		logEntry.FilePath = filePath
	} else if sheetID, ok := metadata["spreadsheetId"].(string); ok {
		logEntry.FilePath = sheetID
	}

	if sheetName, ok := metadata["sheetName"].(string); ok {
		logEntry.SheetName = sheetName
	}

	if changeType, ok := metadata["changeType"].(string); ok {
		logEntry.ChangeType = changeType
	}

	if rowNumber, ok := extractInt(metadata["rowNumber"]); ok {
		logEntry.RowNumber = rowNumber
	}

	if rowData, ok := metadata["rowData"].(map[string]string); ok {
		if rowJSON, err := json.Marshal(rowData); err == nil {
			logEntry.RowData = datatypes.JSON(rowJSON)
		}
	}

	if err := database.DB.Create(&logEntry).Error; err != nil {
		log.Printf("Failed to persist Discord log for area %s: %v", area.Name, err)
	}
}

func buildTemplateVars(area *models.Area, metadata map[string]interface{}) map[string]string {
	vars := map[string]string{
		"areaName":             area.Name,
		"triggerService":       area.TriggerService,
		"actionService":        area.ActionService,
		"eventTitle":           "Scheduled Event",
		"eventTime":            time.Now().Format("2006-01-02 15:04:05"),
		"changeType":           "",
		"sheetName":            "",
		"rowNumber":            "",
		"rowData":              "",
		"rowValues":            "",
		"rowJson":              "",
		"spreadsheetUrl":       "",
		"triggerTime":          "",
		"timerName":            "",
		"interval":             "",
		"messageText":          "",
		"chatId":               "",
		"username":             "",
		"firstName":            "",
		"messageId":            "",
		"openaiGeneratedText":  "",
		"trackId":              "",
		"trackName":            "",
		"artistNames":          "",
		"albumName":            "",
		"trackUrl":             "",
		"previewUrl":           "",
		"deviceName":           "",
		"isPlaying":            "",
		"progressMs":           "",
		"durationMs":           "",
		"coverImageUrl":        "",
		"startedAt":            "",
		"tweetId":              "",
		"tweetText":            "",
		"tweetUrl":             "",
		"tweetAuthorUsername":  "",
		"tweetAuthorName":      "",
		"tweetAuthorId":        "",
		"tweetCreatedAt":       "",
		"tweetLikeCount":       "",
		"tweetNewLikes":        "",
		"tweetRetweetCount":    "",
		"tweetNewRetweets":     "",
		"twitterMonitorType":   "",
		"accountUsername":      "",
		"replyToTweetId":       "",
		"tweetConversationId":  "",
		"tweetInReplyToUserId": "",
		"twitterActionMode":    "",
		"retweetPerformedBy":   "",
		"followerId":           "",
		"followerUsername":     "",
		"followerName":         "",
		"followerBio":          "",
		"followerCreatedAt":    "",
	}

	if metadata == nil {
		return vars
	}

	if changeType, ok := metadata["changeType"].(string); ok {
		vars["changeType"] = changeType
	}
	if sheetName, ok := metadata["sheetName"].(string); ok {
		vars["sheetName"] = sheetName
	}
	if rowNumber, ok := extractInt(metadata["rowNumber"]); ok {
		vars["rowNumber"] = strconv.Itoa(rowNumber)
	}
	if spreadsheetURL, ok := metadata["spreadsheetUrl"].(string); ok {
		vars["spreadsheetUrl"] = spreadsheetURL
	}
	if rowValues, ok := metadata["rowValues"].([]string); ok {
		vars["rowValues"] = strings.Join(rowValues, ", ")
	}
	if rowData, ok := metadata["rowData"].(map[string]string); ok {
		vars["rowData"] = formatRowData(rowData)
		if rowJSON, err := json.Marshal(rowData); err == nil {
			vars["rowJson"] = string(rowJSON)
		}
	}

	if triggerTime, ok := metadata["triggerTime"].(string); ok {
		vars["triggerTime"] = triggerTime
		vars["eventTime"] = triggerTime
	}
	if timerName, ok := metadata["timerName"].(string); ok {
		vars["timerName"] = timerName
	}
	if interval, ok := metadata["interval"].(string); ok {
		vars["interval"] = interval
	}

	if messageText, ok := metadata["messageText"].(string); ok {
		vars["messageText"] = messageText
	}
	if chatID, ok := metadata["chatId"].(string); ok {
		vars["chatId"] = chatID
	}
	if username, ok := metadata["username"].(string); ok {
		vars["username"] = username
	}
	if firstName, ok := metadata["firstName"].(string); ok {
		vars["firstName"] = firstName
	}
	if messageID, ok := extractInt(metadata["messageId"]); ok {
		vars["messageId"] = strconv.Itoa(messageID)
	}
	if trackID, ok := metadata["trackId"].(string); ok {
		vars["trackId"] = trackID
	}
	if trackName, ok := metadata["trackName"].(string); ok {
		vars["trackName"] = trackName
		vars["eventTitle"] = trackName
	}
	if artistNames, ok := metadata["artistNames"].(string); ok {
		vars["artistNames"] = artistNames
	}
	if albumName, ok := metadata["albumName"].(string); ok {
		vars["albumName"] = albumName
	}
	if trackURL, ok := metadata["trackUrl"].(string); ok {
		vars["trackUrl"] = trackURL
	}
	if previewURL, ok := metadata["previewUrl"].(string); ok {
		vars["previewUrl"] = previewURL
	}
	if deviceName, ok := metadata["deviceName"].(string); ok {
		vars["deviceName"] = deviceName
	}
	if coverURL, ok := metadata["coverImageUrl"].(string); ok {
		vars["coverImageUrl"] = coverURL
	}
	if startedAt, ok := metadata["startedAt"].(string); ok {
		vars["startedAt"] = startedAt
		vars["eventTime"] = startedAt
	}
	if isPlayingBool, ok := metadata["isPlaying"].(bool); ok {
		vars["isPlaying"] = strconv.FormatBool(isPlayingBool)
	} else if isPlayingStr, ok := metadata["isPlaying"].(string); ok {
		vars["isPlaying"] = isPlayingStr
	}
	if progressMs, ok := extractInt(metadata["progressMs"]); ok {
		vars["progressMs"] = strconv.Itoa(progressMs)
	}
	if durationMs, ok := extractInt(metadata["durationMs"]); ok {
		vars["durationMs"] = strconv.Itoa(durationMs)
	}

	if openaiText, ok := metadata["openaiGeneratedText"].(string); ok {
		vars["openaiGeneratedText"] = openaiText
	}

	if tweetID, ok := metadata["tweetId"].(string); ok {
		vars["tweetId"] = tweetID
	}
	if tweetText, ok := metadata["tweetText"].(string); ok {
		vars["tweetText"] = tweetText
	}
	if tweetURL, ok := metadata["tweetUrl"].(string); ok {
		vars["tweetUrl"] = tweetURL
	}
	if authorUsername, ok := metadata["tweetAuthorUsername"].(string); ok {
		vars["tweetAuthorUsername"] = authorUsername
	}
	if authorName, ok := metadata["tweetAuthorName"].(string); ok {
		vars["tweetAuthorName"] = authorName
	}
	if authorID, ok := metadata["tweetAuthorId"].(string); ok {
		vars["tweetAuthorId"] = authorID
	}
	if createdAt, ok := metadata["tweetCreatedAt"].(string); ok {
		vars["tweetCreatedAt"] = createdAt
		vars["eventTime"] = createdAt
	}
	if likeCount, ok := extractInt(metadata["tweetLikeCount"]); ok {
		vars["tweetLikeCount"] = strconv.Itoa(likeCount)
	}
	if newLikes, ok := extractInt(metadata["tweetNewLikes"]); ok {
		vars["tweetNewLikes"] = strconv.Itoa(newLikes)
	}
	if retweetCount, ok := extractInt(metadata["tweetRetweetCount"]); ok {
		vars["tweetRetweetCount"] = strconv.Itoa(retweetCount)
	}
	if newRetweets, ok := extractInt(metadata["tweetNewRetweets"]); ok {
		vars["tweetNewRetweets"] = strconv.Itoa(newRetweets)
	}
	if monitorType, ok := metadata["twitterMonitorType"].(string); ok {
		vars["twitterMonitorType"] = monitorType
	}
	if accountUsername, ok := metadata["accountUsername"].(string); ok {
		vars["accountUsername"] = accountUsername
	}
	if actionMode, ok := metadata["twitterActionMode"].(string); ok {
		vars["twitterActionMode"] = actionMode
	}
	if retweetBy, ok := metadata["retweetPerformedBy"].(string); ok {
		vars["retweetPerformedBy"] = retweetBy
	}
	if followerID, ok := metadata["followerId"].(string); ok {
		vars["followerId"] = followerID
	}
	if followerUsername, ok := metadata["followerUsername"].(string); ok {
		vars["followerUsername"] = followerUsername
	}
	if followerName, ok := metadata["followerName"].(string); ok {
		vars["followerName"] = followerName
		if followerName != "" {
			vars["eventTitle"] = followerName
		}
	}
	if followerBio, ok := metadata["followerBio"].(string); ok {
		vars["followerBio"] = followerBio
	}
	if followerCreatedAt, ok := metadata["followerCreatedAt"].(string); ok {
		vars["followerCreatedAt"] = followerCreatedAt
	}
	if replyTo, ok := metadata["replyToTweetId"].(string); ok {
		vars["replyToTweetId"] = replyTo
	}
	if convoID, ok := metadata["tweetConversationId"].(string); ok {
		vars["tweetConversationId"] = convoID
	}
	if replyUserID, ok := metadata["tweetInReplyToUserId"].(string); ok {
		vars["tweetInReplyToUserId"] = replyUserID
	}

	return vars
}

func applyTemplateVariables(input string, vars map[string]string) string {
	result := input
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

func splitSheetHeader(rows [][]string, hasHeader bool) ([]string, [][]string) {
	if len(rows) == 0 {
		return nil, nil
	}

	if !hasHeader {
		return nil, copyRows(rows)
	}

	header := make([]string, len(rows[0]))
	copy(header, rows[0])

	data := [][]string{}
	for _, row := range rows[1:] {
		rowCopy := make([]string, len(row))
		copy(rowCopy, row)
		data = append(data, rowCopy)
	}

	return header, data
}

func detectSheetChange(previous, current [][]string) (*sheetRowChange, bool) {
	minLen := len(previous)
	if len(current) < minLen {
		minLen = len(current)
	}

	for i := 0; i < minLen; i++ {
		if !equalStringSlices(previous[i], current[i]) {
			return &sheetRowChange{
				ChangeType: "updated",
				RowNumber:  i + 1,
				Values:     current[i],
			}, true
		}
	}

	if len(current) > len(previous) {
		return &sheetRowChange{
			ChangeType: "added",
			RowNumber:  len(current),
			Values:     current[len(current)-1],
		}, true
	}

	if len(current) < len(previous) {
		return &sheetRowChange{
			ChangeType: "removed",
			RowNumber:  len(previous),
			Values:     previous[len(previous)-1],
		}, true
	}

	return nil, false
}

func hashSheetRows(rows [][]string) string {
	hasher := sha256.New()
	for _, row := range rows {
		hasher.Write([]byte(strings.Join(row, "|")))
		hasher.Write([]byte("\n"))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func buildRowMap(headers []string, values []string) map[string]string {
	rowData := make(map[string]string)

	for idx, value := range values {
		key := fmt.Sprintf("column_%d", idx+1)
		if headers != nil && idx < len(headers) && headers[idx] != "" {
			key = headers[idx]
		}
		rowData[key] = value
	}

	return rowData
}

func (s *SchedulerService) persistTwitterConfig(area models.Area, cfg twitterTriggerConfig) error {
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return database.DB.Model(&area).Update("trigger_config", datatypes.JSON(cfgBytes)).Error
}

func (s *SchedulerService) persistGoogleSheetsConfig(area models.Area, cfg googleSheetsTriggerConfig) error {
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return database.DB.Model(&area).Update("trigger_config", datatypes.JSON(cfgBytes)).Error
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func copyRows(rows [][]string) [][]string {
	copied := make([][]string, len(rows))
	for i, row := range rows {
		rowCopy := make([]string, len(row))
		copy(rowCopy, row)
		copied[i] = rowCopy
	}
	return copied
}

func formatRowData(row map[string]string) string {
	if len(row) == 0 {
		return ""
	}

	keys := make([]string, 0, len(row))
	for key := range row {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", key, row[key]))
	}

	return strings.Join(parts, ", ")
}

func buildTweetURL(username, tweetID string) string {
	username = strings.TrimSpace(username)
	tweetID = strings.TrimSpace(tweetID)
	if username == "" || tweetID == "" {
		return ""
	}
	return fmt.Sprintf("https://twitter.com/%s/status/%s", username, tweetID)
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func extractInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case float32:
		return int(v), true
	case string:
		if v == "" {
			return 0, false
		}
		num, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return num, true
	default:
		return 0, false
	}
}

func resolveColumnIndex(headers []string, identifier string) int {
	normalized := strings.TrimSpace(strings.ToLower(identifier))
	if normalized == "" {
		return -1
	}

	for idx, header := range headers {
		if strings.ToLower(strings.TrimSpace(header)) == normalized {
			return idx
		}
	}

	if idx, err := parseColumnIdentifier(identifier); err == nil {
		return idx
	}

	return -1
}

func parseColumnIdentifier(identifier string) (int, error) {
	value := strings.TrimSpace(identifier)
	if value == "" {
		return -1, fmt.Errorf("empty column identifier")
	}

	if idx, err := strconv.Atoi(value); err == nil {
		if idx <= 0 {
			return -1, fmt.Errorf("column index must be positive")
		}
		return idx - 1, nil
	}

	value = strings.ToUpper(value)
	total := 0
	for _, r := range value {
		if r < 'A' || r > 'Z' {
			return -1, fmt.Errorf("invalid column identifier: %s", identifier)
		}
		total = total*26 + int(r-'A'+1)
	}

	return total - 1, nil
}

func (s *SchedulerService) StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Scheduler started - checking for scheduled areas every 30 seconds")

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			if err := s.CheckScheduledAreas(); err != nil {
				log.Printf("Error checking scheduled areas: %v", err)
			}
		}
	}
}

func (s *SchedulerService) TestScheduler(areaID string) error {
	var area models.Area
	if err := database.DB.Where("id = ?", areaID).First(&area).Error; err != nil {
		return fmt.Errorf("area not found: %v", err)
	}

	var triggerConfig map[string]interface{}
	if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
		return fmt.Errorf("failed to parse trigger config: %v", err)
	}
	triggerConfig["eventTime"] = time.Now().Format(time.RFC3339)

	updatedConfig, _ := json.Marshal(triggerConfig)
	area.TriggerConfig = updatedConfig
	database.DB.Save(&area)

	return s.executeArea(area, nil)
}
