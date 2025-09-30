package database

import (
	"Golang-API-tutoriel/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		{Email: "john@example.com", Password: string(hashedPassword), FirstName: "John", LastName: "Doe"},
		{Email: "jane@example.com", Password: string(hashedPassword), FirstName: "Jane", LastName: "Smith"},
	}

	for _, user := range users {
		var existingUser models.User
		if err := DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			DB.Create(&user)
		}
	}
	services := []models.Service{
		{Name: "Gmail", Description: "Service de messagerie Google", IconURL: "https://gmail.com/icon.png", IsActive: true},
		{Name: "Slack", Description: "Plateforme de communication", IconURL: "https://slack.com/icon.png", IsActive: true},
		{Name: "GitHub", Description: "Plateforme de développement", IconURL: "https://github.com/icon.png", IsActive: true},
		{Name: "Weather", Description: "Service météorologique", IconURL: "https://weather.com/icon.png", IsActive: true},
		{Name: "Google Calendar", Description: "Service de calendrier Google", IconURL: "https://calendar.google.com/icon.png", IsActive: true},
		{Name: "Discord", Description: "Plateforme de communication", IconURL: "https://discord.com/icon.png", IsActive: true},
	}

	for _, service := range services {
		var existingService models.Service
		if err := DB.Where("name = ?", service.Name).First(&existingService).Error; err != nil {
			DB.Create(&service)
		}
	}

	var gmail, slack, github, weather, calendar, discord models.Service
	DB.Where("name = ?", "Gmail").First(&gmail)
	DB.Where("name = ?", "Slack").First(&slack)
	DB.Where("name = ?", "GitHub").First(&github)
	DB.Where("name = ?", "Weather").First(&weather)
	DB.Where("name = ?", "Google Calendar").First(&calendar)
	DB.Where("name = ?", "Discord").First(&discord)

	actions := []models.Action{
		{ServiceID: gmail.ID, Name: "Nouveau email reçu", Description: "Se déclenche quand un nouvel email arrive", Parameters: `{"sender": "", "subject": ""}`},
		{ServiceID: github.ID, Name: "Nouveau commit", Description: "Se déclenche lors d'un nouveau commit", Parameters: `{"repository": "", "branch": ""}`},
		{ServiceID: weather.ID, Name: "Température élevée", Description: "Se déclenche si température > seuil", Parameters: `{"city": "", "temperature": 30}`},
		{ServiceID: calendar.ID, Name: "Nouvel événement", Description: "Se déclenche quand un nouvel événement est créé", Parameters: `{"calendar": "", "event": ""}`},
		{ServiceID: github.ID, Name: "Nouvelle issue", Description: "Se déclenche quand une nouvelle issue est créée", Parameters: `{"repository": "", "issue": ""}`},
	}

	for _, action := range actions {
		var existingAction models.Action
		if err := DB.Where("service_id = ? AND name = ?", action.ServiceID, action.Name).First(&existingAction).Error; err != nil {
			DB.Create(&action)
		}
	}

	reactions := []models.Reaction{
		{ServiceID: slack.ID, Name: "Envoyer message", Description: "Envoie un message sur Slack", Parameters: `{"channel": "", "message": ""}`},
		{ServiceID: gmail.ID, Name: "Envoyer email", Description: "Envoie un email", Parameters: `{"to": "", "subject": "", "body": ""}`},
		{ServiceID: github.ID, Name: "Créer issue", Description: "Crée une nouvelle issue", Parameters: `{"repository": "", "title": "", "body": ""}`},
		{ServiceID: gmail.ID, Name: "Envoyer email de notification", Description: "Envoie un email de notification", Parameters: `{"to": "", "subject": "", "body": ""}`},
		{ServiceID: discord.ID, Name: "Envoyer message Discord", Description: "Envoie un message sur Discord", Parameters: `{"channel": "", "message": ""}`},
	}

	for _, reaction := range reactions {
		var existingReaction models.Reaction
		if err := DB.Where("service_id = ? AND name = ?", reaction.ServiceID, reaction.Name).First(&existingReaction).Error; err != nil {
			DB.Create(&reaction)
		}
	}


	log.Println("Données de test créées avec succès!")
}
