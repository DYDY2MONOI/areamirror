package main

import (
	"Golang-API-tutoriel/controllers"
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("need .env file")
	}

	database.InitDB()

	database.DB.AutoMigrate(&models.User{}, &models.Service{}, &models.Action{}, &models.Reaction{}, &models.Area{}, &models.Role{}, &models.UserRole{}, &models.RefreshToken{}, &models.OAuth2Token{}, &models.DiscordMessageLog{})

	database.SeedData()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/about.json", controllers.AboutJSON)

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/profile", controllers.AuthMiddleware(), controllers.GetProfile)

	r.POST("/oauth2/login", controllers.OAuth2Login)
	r.POST("/oauth2/refresh", controllers.RefreshToken)
	r.GET("/oauth2/me", controllers.AuthMiddleware(), controllers.GetMe)

	r.GET("/oauth2/github/callback", controllers.GitHubDirectLogin)
	r.GET("/oauth2/google/callback", controllers.GoogleDirectLogin)
	r.GET("/oauth2/spotify/callback", controllers.SpotifyDirectLogin)
	r.GET("/oauth2/facebook/callback", controllers.FacebookDirectLogin)
	r.GET("/oauth2/twitter/callback", controllers.TwitterDirectLogin)

	r.POST("/mobile/oauth2/login", controllers.MobileOAuth2Login)
	r.POST("/mobile/oauth2/refresh", controllers.RefreshToken)
	r.GET("/mobile/oauth2/me", controllers.AuthMiddleware(), controllers.GetMe)
	r.GET("/mobile/user/me/areas", controllers.AuthMiddleware(), controllers.GetUserAreas)
	r.GET("/mobile/areas/popular", controllers.GetPopularAreas)
	r.GET("/mobile/areas/recommended", controllers.GetRecommendedAreas)
	r.POST("/mobile/areas", controllers.AuthMiddleware(), controllers.CreateArea)
	r.PUT("/mobile/areas/:id", controllers.AuthMiddleware(), controllers.UpdateArea)
	r.DELETE("/mobile/areas/:id", controllers.AuthMiddleware(), controllers.DeleteArea)
	r.PATCH("/mobile/areas/:id/toggle", controllers.AuthMiddleware(), controllers.ToggleArea)

	r.PUT("/profile", controllers.AuthMiddleware(), controllers.UpdateProfile)
	r.POST("/profile/image", controllers.AuthMiddleware(), controllers.UploadProfileImage)
	r.POST("/profile/github/link", controllers.AuthMiddleware(), controllers.LinkGitHubAccount)
	r.DELETE("/profile/github/unlink", controllers.AuthMiddleware(), controllers.UnlinkGitHubAccount)
	r.POST("/profile/google/link", controllers.AuthMiddleware(), controllers.LinkGoogleAccount)
	r.DELETE("/profile/google/unlink", controllers.AuthMiddleware(), controllers.UnlinkGoogleAccount)
	r.POST("/profile/facebook/link", controllers.AuthMiddleware(), controllers.LinkFacebookAccount)
	r.DELETE("/profile/facebook/unlink", controllers.AuthMiddleware(), controllers.UnlinkFacebookAccount)
	r.POST("/profile/onedrive/link", controllers.AuthMiddleware(), controllers.LinkOneDriveAccount)
	r.DELETE("/profile/onedrive/unlink", controllers.AuthMiddleware(), controllers.UnlinkOneDriveAccount)
	r.POST("/profile/spotify/link", controllers.AuthMiddleware(), controllers.LinkSpotifyAccount)
	r.DELETE("/profile/spotify/unlink", controllers.AuthMiddleware(), controllers.UnlinkSpotifyAccount)
	r.POST("/profile/twitter/link", controllers.AuthMiddleware(), controllers.LinkTwitterAccount)
	r.DELETE("/profile/twitter/unlink", controllers.AuthMiddleware(), controllers.UnlinkTwitterAccount)

	r.GET("/gmail/oauth2/setup", controllers.AuthMiddleware(), controllers.SetupGmailOAuth2)
	r.POST("/gmail/oauth2/token", controllers.AuthMiddleware(), controllers.StoreGmailToken)
	r.GET("/gmail/oauth2/status", controllers.AuthMiddleware(), controllers.GetGmailTokenStatus)
	r.POST("/gmail/oauth2/test", controllers.AuthMiddleware(), controllers.TestGmailConnection)
	r.DELETE("/gmail/oauth2/revoke", controllers.AuthMiddleware(), controllers.RevokeGmailToken)

	// OneDrive routes
	r.GET("/onedrive/auth/start", controllers.OneDriveAuthStart)
	r.GET("/onedrive/callback", controllers.OneDriveCallback)
	r.GET("/onedrive/files", controllers.OneDriveListFiles)
	r.POST("/onedrive/upload", controllers.OneDriveUploadFile)
	r.GET("/onedrive/download/:fileId", controllers.OneDriveDownloadFile)
	r.DELETE("/onedrive/delete/:fileId", controllers.OneDriveDeleteFile)
	r.POST("/onedrive/folder", controllers.OneDriveCreateFolder)
	r.GET("/onedrive/user", controllers.OneDriveUserInfo)

	// Google Agenda routes
	googleAgendaController := controllers.NewGoogleAgendaController()
	r.GET("/google-agenda/auth", controllers.AuthMiddleware(), googleAgendaController.GetAuthURL)
	r.GET("/google-agenda/callback", controllers.AuthMiddleware(), googleAgendaController.HandleCallback)
	r.GET("/google-agenda/events", controllers.AuthMiddleware(), googleAgendaController.GetUpcomingEvents)
	r.GET("/google-agenda/test", controllers.AuthMiddleware(), googleAgendaController.TestAgendaConnection)
	r.GET("/google-agenda/calendars", controllers.AuthMiddleware(), googleAgendaController.ListCalendars)

	api := r.Group("/api")
	{
		api.GET("/github/repositories", controllers.AuthMiddleware(), controllers.GetGitHubRepositories)
		api.POST("/areas/github-gmail", controllers.AuthMiddleware(), controllers.CreateGitHubGmailArea)
	}

	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUser)
	r.POST("/users", controllers.CreateUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	r.GET("/services", controllers.GetServices)
	r.GET("/services/:id", controllers.GetService)
	r.POST("/services", controllers.CreateService)
	r.PUT("/services/:id", controllers.UpdateService)
	r.DELETE("/services/:id", controllers.DeleteService)

	r.GET("/actions", controllers.GetActions)
	r.GET("/actions/:id", controllers.GetAction)
	r.POST("/actions", controllers.CreateAction)
	r.PUT("/actions/:id", controllers.UpdateAction)
	r.DELETE("/actions/:id", controllers.DeleteAction)

	r.GET("/reactions", controllers.GetReactions)
	r.GET("/reactions/:id", controllers.GetReaction)
	r.POST("/reactions", controllers.CreateReaction)
	r.PUT("/reactions/:id", controllers.UpdateReaction)
	r.DELETE("/reactions/:id", controllers.DeleteReaction)

	r.GET("/service/:id/actions", controllers.GetServiceActions)
	r.GET("/service/:id/reactions", controllers.GetServiceReactions)

	r.GET("/areas", controllers.AuthMiddleware(), controllers.GetAreas)
	r.GET("/areas/:id", controllers.GetArea)
	r.GET("/areas/:id/discord-logs", controllers.AuthMiddleware(), controllers.GetAreaDiscordLogs)
	r.GET("/test-area/:id", controllers.GetArea)
	r.POST("/areas", controllers.AuthMiddleware(), controllers.CreateArea)
	r.PUT("/areas/:id", controllers.AuthMiddleware(), controllers.UpdateArea)
	r.DELETE("/areas/:id", controllers.AuthMiddleware(), controllers.DeleteArea)
	r.PATCH("/areas/:id/toggle", controllers.AuthMiddleware(), controllers.ToggleArea)

	r.GET("/user/me/areas", controllers.AuthMiddleware(), controllers.GetUserAreas)

	r.POST("/user/:id/applets", controllers.CreateApplet)
	r.GET("/user/:id/applets", controllers.GetApplets)
	r.GET("/user/:id/applets/:id", controllers.GetApplet)
	r.PUT("/user/:id/applets/:id", controllers.UpdateApplet)
	r.DELETE("/user/:id/applets/:id", controllers.DeleteApplet)

	githubWebhookController := controllers.NewGitHubWebhookController()
	r.POST("/webhooks/github", githubWebhookController.HandleWebhook)
	r.POST("/webhooks/telegram", controllers.TelegramWebhook)

	r.Static("/uploads", "./uploads")
	r.StaticFile("/test-onedrive.html", "./test-onedrive.html")
	r.StaticFile("/test-slack.html", "./test-slack.html")
	r.StaticFile("/test-slack-advanced.html", "./test-slack-advanced.html")

	r.GET("/areas/popular", controllers.GetPopularAreas)
	r.GET("/areas/recommended", controllers.GetRecommendedAreas)

	r.POST("/test/email", controllers.TestEmail)
	r.POST("/test/discord", controllers.TestDiscord)
	r.POST("/test/slack", controllers.TestSlack)
	r.GET("/test/spotify", controllers.AuthMiddleware(), controllers.TestSpotify)
	r.POST("/test/google-sheets", controllers.TestGoogleSheets)
	r.POST("/test/weather", controllers.TestWeatherTrigger)
	r.GET("/weather", controllers.GetWeatherData)
	r.POST("/test/scheduler/:id", controllers.TestScheduler)

	r.GET("/roles", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.GetRoles)
	r.POST("/roles", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.CreateRole)
	r.GET("/roles/:id", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.GetRole)
	r.PUT("/roles/:id", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.UpdateRole)
	r.DELETE("/roles/:id", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.DeleteRole)

	r.POST("/users/:id/roles", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.AssignRoleToUser)
	r.DELETE("/users/:id/roles", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.RemoveRoleFromUser)
	r.GET("/users/:id/roles", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.GetUserRoles)
	r.PUT("/users/:id/role", controllers.AuthMiddleware(), controllers.RoleMiddleware("admin"), controllers.UpdateUserRole)

	scheduler, err := services.NewSchedulerService()
	if err != nil {
		log.Printf("Warning: Failed to initialize scheduler: %v", err)
	} else {
		ctx := context.Background()
		go scheduler.StartScheduler(ctx)
		log.Println("Scheduler started in background")
	}

	// Setup Telegram webhook if configured
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	webhookURL := os.Getenv("TELEGRAM_WEBHOOK_URL")
	if botToken != "" && webhookURL != "" {
		if err := controllers.SetupTelegramWebhook(botToken, webhookURL); err != nil {
			log.Printf("Warning: Failed to setup Telegram webhook: %v", err)
		}
	}

	r.Run()
}
