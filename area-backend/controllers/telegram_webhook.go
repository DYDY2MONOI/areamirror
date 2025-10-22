package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TelegramWebhook(c *gin.Context) {
	var update services.TelegramUpdate
	if err := c.BindJSON(&update); err != nil {
		log.Printf("Failed to parse Telegram webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("Received Telegram webhook: Update ID %d, Chat ID %d, Text: %s",
		update.UpdateID, update.Message.Chat.ID, update.Message.Text)

	chatIDStr := strconv.FormatInt(update.Message.Chat.ID, 10)

	go handleTelegramTriggers(update, chatIDStr)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleTelegramTriggers(update services.TelegramUpdate, chatID string) {
	var areas []models.Area

	err := database.DB.Where("trigger_service = ? AND is_active = ?", "Telegram", true).Find(&areas).Error
	if err != nil {
		log.Printf("Failed to fetch Telegram trigger areas: %v", err)
		return
	}

	log.Printf("Found %d active Telegram trigger areas", len(areas))

	for _, area := range areas {
		var triggerConfig map[string]interface{}
		if err := json.Unmarshal(area.TriggerConfig, &triggerConfig); err != nil {
			log.Printf("Failed to parse trigger config for area %s: %v", area.Name, err)
			continue
		}

		expectedChatID, _ := triggerConfig["chatId"].(string)
		if expectedChatID == "" {
			expectedChatID, _ = triggerConfig["chatID"].(string)
		}
		if expectedChatID == "" {
			expectedChatID, _ = triggerConfig["chat_id"].(string)
		}

		if expectedChatID != "" && expectedChatID != chatID {
			log.Printf("Chat ID mismatch for area %s: expected %s, got %s", area.Name, expectedChatID, chatID)
			continue
		}

		triggerType := area.TriggerType
		messageText := update.Message.Text

		shouldTrigger := false
		switch triggerType {
		case "message_received":
			shouldTrigger = true
		case "keyword_match":
			keyword, _ := triggerConfig["keyword"].(string)
			if keyword != "" && containsKeyword(messageText, keyword) {
				shouldTrigger = true
			}
		case "command_received":
			command, _ := triggerConfig["command"].(string)
			if command != "" && messageText == command {
				shouldTrigger = true
			}
		default:
			log.Printf("Unknown trigger type %s for area %s", triggerType, area.Name)
		}

		if shouldTrigger {
			log.Printf("Triggering area %s from Telegram message", area.Name)

			metadata := map[string]interface{}{
				"messageText": messageText,
				"chatId":      chatID,
				"username":    update.Message.From.Username,
				"firstName":   update.Message.From.FirstName,
				"messageId":   update.Message.MessageID,
			}

			scheduler, err := services.NewSchedulerService()
			if err != nil {
				log.Printf("Failed to create scheduler service: %v", err)
				continue
			}

			if err := scheduler.ExecuteAreaPublic(area, metadata); err != nil {
				log.Printf("Failed to execute area %s: %v", area.Name, err)
			}
		}
	}
}

func containsKeyword(text, keyword string) bool {
	if keyword == "" {
		return false
	}
	return containsSubstring(text, keyword)
}

func containsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func SetupTelegramWebhook(botToken, webhookURL string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", botToken, webhookURL)
	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to setup webhook: %w", err)
	}
	defer resp.Body.Close()

	var result services.TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.OK {
		return fmt.Errorf("telegram API error: %s", result.Description)
	}

	log.Printf("Telegram webhook setup successfully: %s", webhookURL)
	return nil
}

