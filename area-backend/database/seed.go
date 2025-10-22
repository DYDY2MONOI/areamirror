package database

import (
	"Golang-API-tutoriel/models"
	"encoding/json"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	createDefaultRoles()

	users := []models.User{
		{Email: "admin@area.com", Password: string(adminPassword), FirstName: "Admin", LastName: "User", Role: "admin"},
		{Email: "john@example.com", Password: string(hashedPassword), FirstName: "John", LastName: "Doe", Role: "member"},
		{Email: "jane@example.com", Password: string(hashedPassword), FirstName: "Jane", LastName: "Smith", Role: "member"},
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
		{Name: "OneDrive", Description: "Service de stockage Microsoft", IconURL: "https://onedrive.com/icon.png", IsActive: true},
	}

	for _, service := range services {
		var existingService models.Service
		if err := DB.Where("name = ?", service.Name).First(&existingService).Error; err != nil {
			DB.Create(&service)
		}
	}

	var gmail, slack, github, weather, calendar, discord, onedrive models.Service
	DB.Where("name = ?", "Gmail").First(&gmail)
	DB.Where("name = ?", "Slack").First(&slack)
	DB.Where("name = ?", "GitHub").First(&github)
	DB.Where("name = ?", "Weather").First(&weather)
	DB.Where("name = ?", "Google Calendar").First(&calendar)
	DB.Where("name = ?", "Discord").First(&discord)
	DB.Where("name = ?", "OneDrive").First(&onedrive)

	actions := []models.Action{
		{ServiceID: gmail.ID, Name: "Nouveau email reçu", Description: "Se déclenche quand un nouvel email arrive", Parameters: `{"sender": "", "subject": ""}`},
		{ServiceID: github.ID, Name: "Nouveau commit", Description: "Se déclenche lors d'un nouveau commit", Parameters: `{"repository": "", "branch": ""}`},
		{ServiceID: weather.ID, Name: "Température élevée", Description: "Se déclenche si température > seuil", Parameters: `{"city": "", "temperature": 30}`},
		{ServiceID: calendar.ID, Name: "Nouvel événement", Description: "Se déclenche quand un nouvel événement est créé", Parameters: `{"calendar": "", "event": ""}`},
		{ServiceID: github.ID, Name: "Nouvelle issue", Description: "Se déclenche quand une nouvelle issue est créée", Parameters: `{"repository": "", "issue": ""}`},
		{ServiceID: onedrive.ID, Name: "Nouveau fichier", Description: "Se déclenche quand un nouveau fichier est ajouté", Parameters: `{"folder": "", "fileName": ""}`},
		{ServiceID: onedrive.ID, Name: "Fichier modifié", Description: "Se déclenche quand un fichier est modifié", Parameters: `{"folder": "", "fileName": ""}`},
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
		{ServiceID: discord.ID, Name: "Envoyer message Discord", Description: "Envoie un message sur Discord", Parameters: `{"webhookUrl": "", "message": ""}`},
		{ServiceID: onedrive.ID, Name: "Upload fichier", Description: "Upload un fichier vers OneDrive", Parameters: `{"fileName": "", "content": ""}`},
		{ServiceID: onedrive.ID, Name: "Créer dossier", Description: "Crée un nouveau dossier sur OneDrive", Parameters: `{"folderName": ""}`},
	}

	for _, reaction := range reactions {
		var existingReaction models.Reaction
		if err := DB.Where("service_id = ? AND name = ?", reaction.ServiceID, reaction.Name).First(&existingReaction).Error; err != nil {
			DB.Create(&reaction)
		}
	}

	log.Println("Données de test créées avec succès!")
}

func createDefaultRoles() {
	adminPermissions := models.GetDefaultPermissions(models.RoleAdmin)
	adminPermissionsJSON, _ := json.Marshal(adminPermissions)

	adminRole := models.Role{
		Name:        models.RoleAdmin,
		Description: "Administrator with full system access",
		Permissions: string(adminPermissionsJSON),
		IsActive:    true,
	}

	var existingAdminRole models.Role
	if err := DB.Where("name = ?", models.RoleAdmin).First(&existingAdminRole).Error; err != nil {
		DB.Create(&adminRole)
	}

	memberPermissions := models.GetDefaultPermissions(models.RoleMember)
	memberPermissionsJSON, _ := json.Marshal(memberPermissions)

	memberRole := models.Role{
		Name:        models.RoleMember,
		Description: "Regular user with basic access",
		Permissions: string(memberPermissionsJSON),
		IsActive:    true,
	}

	var existingMemberRole models.Role
	if err := DB.Where("name = ?", models.RoleMember).First(&existingMemberRole).Error; err != nil {
		DB.Create(&memberRole)
	}

	log.Println("Default roles created successfully!")
}
