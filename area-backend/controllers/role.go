package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range requiredRoles {
			if user.Role == role {
				hasRole = true
				break
			}
			for _, userRole := range user.Roles {
				if userRole.Name == role {
					hasRole = true
					break
				}
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("userRole", user.Role)
		c.Next()
	}
}

func PermissionMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, permission := range requiredPermissions {
			if hasUserPermission(&user, permission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasUserPermission(user *models.User, permission string) bool {
	if user.Role == "admin" {
		return true
	}

	for _, role := range user.Roles {
		var permissions []string
		if err := json.Unmarshal([]byte(role.Permissions), &permissions); err == nil {
			for _, p := range permissions {
				if p == permission {
					return true
				}
			}
		}
	}

	return false
}

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role.Permissions == "" {
		permissions := models.GetDefaultPermissions(role.Name)
		permissionsJSON, _ := json.Marshal(permissions)
		role.Permissions = string(permissionsJSON)
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

func GetRoles(c *gin.Context) {
	var roles []models.Role
	database.DB.Find(&roles)
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func GetRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func UpdateRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var input models.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&role).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": role})
}

func DeleteRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	database.DB.Delete(&role)
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

func AssignRoleToUser(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var role models.Role
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	userRole := models.UserRole{
		UserID: req.UserID,
		RoleID: req.RoleID,
	}

	if err := database.DB.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to assign role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role assigned successfully"})
}

func RemoveRoleFromUser(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userRole models.UserRole
	if err := database.DB.Where("user_id = ? AND role_id = ?", req.UserID, req.RoleID).First(&userRole).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role assignment not found"})
		return
	}

	database.DB.Delete(&userRole)
	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}

func GetUserRoles(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user.Roles})
}

func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validRoles := []string{"admin", "member"}
	if !contains(validRoles, req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = req.Role
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully", "user": user})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
