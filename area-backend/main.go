package main

import (
	"Golang-API-tutoriel/controllers"
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/profile", controllers.AuthMiddleware(), controllers.GetProfile)

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

	r.GET("/user/:id/areas", controllers.GetUserAreas)

	r.Run()
}
