package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your-secret-key-change-in-production")

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ProfileUpdateRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	Country        string `json:"country"`
	CurrentPassword string `json:"current_password"`
	NewPassword    string `json:"new_password"`
}

type GitHubLinkRequest struct {
	Code string `json:"code" binding:"required"`
}

type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GitHubUserResponse struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func init() {
	if key := os.Getenv("JWT_SECRET"); key != "" {
		jwtKey = []byte(key)
	}
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A user with this email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return
	}

	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"profile_image": user.ProfileImage,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}

	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"profile_image": user.ProfileImage,
		},
	})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			"phone":      user.Phone,
			"birthday":   user.Birthday,
			"gender":     user.Gender,
			"country":    user.Country,
			"lang":       user.Lang,
			"login_provider": user.LoginProvider,
			"profile_image": user.ProfileImage,
			"github_id":  user.GitHubID,
			"github_username": user.GitHubUsername,
		},
	})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Country != "" {
		updates["country"] = req.Country
	}

	if req.CurrentPassword != "" && req.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password incorrect"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing new password"})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			"phone":      user.Phone,
			"birthday":   user.Birthday,
			"gender":     user.Gender,
			"country":    user.Country,
			"lang":       user.Lang,
			"login_provider": user.LoginProvider,
			"profile_image": user.ProfileImage,
			"github_id":  user.GitHubID,
			"github_username": user.GitHubUsername,
		},
	})
}

func generateJWT(userID uint, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token required"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Next()
	}
}

func UploadProfileImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image must not exceed 5MB"})
		return
	}

	uploadDir := "uploads/profile_images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating directory"})
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("profile_%d_%d%s", userID, time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ProfileImage != nil && *user.ProfileImage != "" {
		oldPath := *user.ProfileImage
		if strings.HasPrefix(oldPath, "uploads/") {
			os.Remove(oldPath)
		}
	}

	profileImagePath := filepath
	if err := database.DB.Model(&user).Update("profile_image", profileImagePath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			"phone":      user.Phone,
			"birthday":   user.Birthday,
			"gender":     user.Gender,
			"country":    user.Country,
			"lang":       user.Lang,
			"login_provider": user.LoginProvider,
			"profile_image": user.ProfileImage,
			"github_id":  user.GitHubID,
			"github_username": user.GitHubUsername,
		},
	})
}


func LinkGitHubAccount(c *gin.Context) {
userID, exists := c.Get("userID")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
return
}

var req GitHubLinkRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

githubClientID := os.Getenv("GITHUB_CLIENT_ID")
githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

if githubClientID == "" || githubClientSecret == "" {
c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub OAuth not configured"})
return
}

accessToken, err := exchangeCodeForToken(req.Code, githubClientID, githubClientSecret)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
return
}

githubUser, err := getGitHubUser(accessToken)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get GitHub user"})
return
}

var user models.User
if err := database.DB.First(&user, userID).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
return
}

githubIDStr := fmt.Sprintf("%d", githubUser.ID)

var existingUser models.User
if err := database.DB.Where("github_id = ?", githubIDStr).First(&existingUser).Error; err == nil {
c.JSON(http.StatusConflict, gin.H{"error": "This GitHub account is already linked to another user"})
return
}

user.GitHubID = &githubIDStr
user.GitHubUsername = &githubUser.Login

if err := database.DB.Save(&user).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link GitHub account"})
return
}

c.JSON(http.StatusOK, gin.H{
"message": "GitHub account linked successfully",
"github_username": user.GitHubUsername,
})
}

func UnlinkGitHubAccount(c *gin.Context) {
userID, exists := c.Get("userID")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
return
}

var user models.User
if err := database.DB.First(&user, userID).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
return
}

if user.GitHubID == nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "No GitHub account linked"})
return
}

user.GitHubID = nil
user.GitHubUsername = nil

if err := database.DB.Save(&user).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink GitHub account"})
return
}

c.JSON(http.StatusOK, gin.H{
"message": "GitHub account unlinked successfully",
})
}

func exchangeCodeForToken(code, clientID, clientSecret string) (string, error) {
url := "https://github.com/login/oauth/access_token"

data := map[string]string{
"client_id":     clientID,
"client_secret": clientSecret,
"code":          code,
}

jsonData, err := json.Marshal(data)
if err != nil {
return "", err
}

req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
if err != nil {
return "", err
}

req.Header.Set("Accept", "application/json")
req.Header.Set("Content-Type", "application/json")

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}

var tokenResp GitHubTokenResponse
if err := json.Unmarshal(body, &tokenResp); err != nil {
return "", err
}

return tokenResp.AccessToken, nil
}

func getGitHubUser(accessToken string) (*GitHubUserResponse, error) {
url := "https://api.github.com/user"

req, err := http.NewRequest("GET", url, nil)
if err != nil {
return nil, err
}

req.Header.Set("Authorization", "Bearer "+accessToken)
req.Header.Set("Accept", "application/vnd.github.v3+json")

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return nil, err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return nil, err
}

var githubUser GitHubUserResponse
if err := json.Unmarshal(body, &githubUser); err != nil {
return nil, err
}

return &githubUser, nil
}
