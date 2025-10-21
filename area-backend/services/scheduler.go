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
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type SchedulerService struct {
	emailService   *EmailService
	discordService *DiscordService
	weatherService *WeatherService
	sheetsService  *GoogleSheetsService
  driveService   *GoogleDriveService
	telegramService *TelegramService
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

	return &SchedulerService{
		emailService:    emailService,
		discordService:  discordService,
		weatherService:  weatherService,
		sheetsService:   sheetsService,
    driveService:   driveService,
		telegramService: telegramService,
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

	return nil
}

type googleDriveTriggerConfig struct {
    FolderID     string            `json:"folderId"`
    KnownFileIDs map[string]bool   `json:"knownFileIds"`
    LastChecked  *time.Time        `json:"lastChecked"`
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
	if err := s.checkTimerTriggers(); err != nil {
		log.Printf("Error checking timer triggers: %v", err)
	}

	return nil
}

func (s *SchedulerService) checkCalendarTriggers() error {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Google Calendar", true).Find(&areas).Error
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

	if area.LastRunAt != nil {
		timeSinceLastRun := now.Sub(*area.LastRunAt)
		if timeSinceLastRun < 5*time.Minute {
			log.Printf("Area %s already executed recently, skipping", area.Name)
			return false
		}
	}

	timeDiff := eventTime.Sub(now)
	return timeDiff >= 0 && timeDiff <= 30*time.Second
}

func (s *SchedulerService) shouldTriggerWeatherArea(area models.Area, triggerConfig map[string]interface{}) bool {
	if s.weatherService == nil {
		log.Printf("Weather service not available for area %s", area.Name)
		return false
	}

	if area.LastRunAt != nil {
		timeSinceLastRun := time.Since(*area.LastRunAt)
		if timeSinceLastRun < 10*time.Minute {
			log.Printf("Weather area %s already executed recently, skipping", area.Name)
			return false
		}
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

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Printf("Failed to parse interval '%s' for timer area %s: %v", intervalStr, area.Name, err)
		return false
	}

	if area.LastRunAt != nil {
		timeSinceLastRun := now.Sub(*area.LastRunAt)
		if timeSinceLastRun < interval {
			return false
		}
	}

	log.Printf("Timer trigger activated for area %s (interval: %s)", area.Name, intervalStr)
	return true
}

func (s *SchedulerService) executeArea(area models.Area, metadata map[string]interface{}) error {
	log.Printf("Executing area: %s", area.Name)

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
	default:
		log.Printf("Unsupported action service: %s", area.ActionService)
		return nil
	}
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
	message = applyTemplateVariables(message, templateVars)

	if err := s.telegramService.SendMessage(chatID, message); err != nil {
		s.recordAreaFailure(area, fmt.Errorf("failed to send telegram message: %w", err))
		return fmt.Errorf("failed to send telegram message: %v", err)
	}

	s.recordAreaSuccess(area)
	log.Printf("Telegram message sent successfully for AREA: %s", area.Name)
	log.Printf("Successfully executed area: %s", area.Name)
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
		"areaName":       area.Name,
		"triggerService": area.TriggerService,
		"actionService":  area.ActionService,
		"eventTitle":     "Scheduled Event",
		"eventTime":      time.Now().Format("2006-01-02 15:04:05"),
		"changeType":     "",
		"sheetName":      "",
		"rowNumber":      "",
		"rowData":        "",
		"rowValues":      "",
		"rowJson":        "",
		"spreadsheetUrl": "",
		"triggerTime":    "",
		"timerName":      "",
		"interval":       "",
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
