package main

import (
	"Golang-API-tutoriel/controllers"
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"context"
	"log"

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

	database.DB.AutoMigrate(&models.User{}, &models.Service{}, &models.Action{}, &models.Reaction{}, &models.Area{})

	database.SeedData()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes d'authentification directement accessibles (pour compatibilité avec le frontend)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/profile", controllers.AuthMiddleware(), controllers.GetProfile)
	r.PUT("/profile", controllers.AuthMiddleware(), controllers.UpdateProfile)
	r.POST("/profile/image", controllers.AuthMiddleware(), controllers.UploadProfileImage)
	r.POST("/profile/github/link", controllers.AuthMiddleware(), controllers.LinkGitHubAccount)
	r.DELETE("/profile/github/unlink", controllers.AuthMiddleware(), controllers.UnlinkGitHubAccount)
	r.POST("/profile/google/link", controllers.AuthMiddleware(), controllers.LinkGoogleAccount)
	r.DELETE("/profile/google/unlink", controllers.AuthMiddleware(), controllers.UnlinkGoogleAccount)

	// Routes GitHub dans le groupe /api
	api := r.Group("/api")
	{
		api.GET("/github/repositories", controllers.AuthMiddleware(), controllers.GetGitHubRepositories)
		api.POST("/areas/github-gmail", controllers.AuthMiddleware(), controllers.CreateGitHubGmailArea)
	}

	// Autres routes directement accessibles
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

	r.GET("/areas", controllers.GetAreas)
	r.GET("/areas/:id", controllers.GetArea)
	r.POST("/areas", controllers.CreateArea)
	r.PUT("/areas/:id", controllers.UpdateArea)
	r.DELETE("/areas/:id", controllers.DeleteArea)
	r.PATCH("/areas/:id/toggle", controllers.ToggleArea)

	r.GET("/user/me/areas", controllers.AuthMiddleware(), controllers.GetUserAreas)

	r.POST("/user/:id/applets", controllers.CreateApplet)
	r.GET("/user/:id/applets", controllers.GetApplets)
	r.GET("/user/:id/applets/:id", controllers.GetApplet)
	r.PUT("/user/:id/applets/:id", controllers.UpdateApplet)
	r.DELETE("/user/:id/applets/:id", controllers.DeleteApplet)

	githubWebhookController := controllers.NewGitHubWebhookController()
	r.POST("/webhooks/github", githubWebhookController.HandleWebhook)

	r.Static("/uploads", "./uploads")

	r.GET("/areas/popular", controllers.GetPopularAreas)
	r.GET("/areas/recommended", controllers.GetRecommendedAreas)

	r.POST("/test/email", controllers.TestEmail)
	r.POST("/test/scheduler/:id", controllers.TestScheduler)

	scheduler, err := services.NewSchedulerService()
	if err != nil {
		log.Printf("Warning: Failed to initialize scheduler: %v", err)
	} else {
		ctx := context.Background()
		go scheduler.StartScheduler(ctx)
		log.Println("Scheduler started in background")
	}

	r.Run()
}
