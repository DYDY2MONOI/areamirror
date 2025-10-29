package controllers

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"Golang-API-tutoriel/services"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

type OAuth2LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         gin.H  `json:"user"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ProfileUpdateRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Phone           string `json:"phone"`
	Country         string `json:"country"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
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
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GoogleLinkRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type GoogleUserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type FacebookLinkRequest struct {
	Code string `json:"code" binding:"required"`
}

type FacebookTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type FacebookUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SpotifyLinkRequest struct {
	Code string `json:"code" binding:"required"`
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyUserProfile struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type TwitterLinkRequest struct {
	Code         string `json:"code" binding:"required"`
	CodeVerifier string `json:"code_verifier" binding:"required"`
}

type TwitterTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type TwitterUserResponse struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}

func init() {
	if key := os.Getenv("JWT_SECRET"); key != "" {
		jwtKey = []byte(key)
	}
}

func getRedirectURI(defaultEnv, mobileEnv, fallback string, isMobile bool) string {
	if isMobile {
		if mobile := os.Getenv(mobileEnv); mobile != "" {
			return mobile
		}
	}

	if standard := os.Getenv(defaultEnv); standard != "" {
		return standard
	}

	return fallback
}

type mobileState struct {
	isMobile bool
	mode     string
	payload  string
}

func parseMobileState(state string) mobileState {
	result := mobileState{mode: "login"}
	if state == "" {
		return result
	}

	if !strings.HasPrefix(state, "mobile") {
		return result
	}

	result.isMobile = true
	parts := strings.Split(state, ":")
	if len(parts) >= 2 {
		if parts[1] == "link" {
			result.mode = "link"
			if len(parts) >= 3 {
				result.payload = parts[2]
			}
		} else if parts[1] != "" {
			result.payload = parts[1]
		}
	}

	if len(parts) >= 3 && parts[1] != "link" {
		result.payload = parts[2]
	}

	return result
}

func redirectToMobileWithTokens(c *gin.Context, provider string, tokenResponse OAuth2TokenResponse) {
	scheme := os.Getenv("MOBILE_CALLBACK_SCHEME")
	if scheme == "" {
		scheme = "area"
	}

	host := os.Getenv("MOBILE_CALLBACK_HOST")
	if host == "" {
		host = "oauth2"
	}

	path := os.Getenv("MOBILE_CALLBACK_PATH")
	if path == "" {
		path = "/callback"
	}

	values := url.Values{}
	values.Set("mode", "login")
	values.Set("provider", provider)
	values.Set("access_token", tokenResponse.AccessToken)
	if tokenResponse.RefreshToken != "" {
		values.Set("refresh_token", tokenResponse.RefreshToken)
	}
	values.Set("token_type", tokenResponse.TokenType)
	expiresIn := tokenResponse.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 900
	}
	values.Set("expires_in", strconv.Itoa(expiresIn))

	redirectURL := fmt.Sprintf("%s://%s%s?%s", scheme, host, path, values.Encode())
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func redirectToMobileWithAuthorizationCode(c *gin.Context, provider, code, codeVerifier string) {
	scheme := os.Getenv("MOBILE_CALLBACK_SCHEME")
	if scheme == "" {
		scheme = "area"
	}

	host := os.Getenv("MOBILE_CALLBACK_HOST")
	if host == "" {
		host = "oauth2"
	}

	path := os.Getenv("MOBILE_CALLBACK_PATH")
	if path == "" {
		path = "/callback"
	}

	values := url.Values{}
	values.Set("mode", "link")
	values.Set("provider", provider)
	values.Set("code", code)
	if codeVerifier != "" {
		values.Set("code_verifier", codeVerifier)
	}

	redirectURL := fmt.Sprintf("%s://%s%s?%s", scheme, host, path, values.Encode())
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
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
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
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
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
		},
	})
}

func OAuth2Login(c *gin.Context) {
	var req OAuth2LoginRequest

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

	accessToken, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	c.JSON(http.StatusOK, OAuth2TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: gin.H{
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
		},
	})
}

func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refreshToken models.RefreshToken
	if err := database.DB.Where("token = ?", req.RefreshToken).First(&refreshToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if !refreshToken.IsValid() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired or revoked"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, refreshToken.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	accessToken, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   900,
	})
}

func GetMe(c *gin.Context) {
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
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"created_at":       user.CreatedAt,
			"updated_at":       user.UpdatedAt,
			"phone":            user.Phone,
			"birthday":         user.Birthday,
			"gender":           user.Gender,
			"country":          user.Country,
			"lang":             user.Lang,
			"login_provider":   user.LoginProvider,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"github_id":        user.GitHubID,
			"github_username":  user.GitHubUsername,
			"google_id":        user.GoogleID,
			"google_email":     user.GoogleEmail,
			"facebook_id":      user.FacebookID,
			"facebook_email":   user.FacebookEmail,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
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
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"created_at":       user.CreatedAt,
			"updated_at":       user.UpdatedAt,
			"phone":            user.Phone,
			"birthday":         user.Birthday,
			"gender":           user.Gender,
			"country":          user.Country,
			"lang":             user.Lang,
			"login_provider":   user.LoginProvider,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"github_id":        user.GitHubID,
			"github_username":  user.GitHubUsername,
			"google_id":        user.GoogleID,
			"google_email":     user.GoogleEmail,
			"facebook_id":      user.FacebookID,
			"facebook_email":   user.FacebookEmail,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
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
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"created_at":       user.CreatedAt,
			"updated_at":       user.UpdatedAt,
			"phone":            user.Phone,
			"birthday":         user.Birthday,
			"gender":           user.Gender,
			"country":          user.Country,
			"lang":             user.Lang,
			"login_provider":   user.LoginProvider,
			"profile_image":    user.ProfileImage,
			"github_id":        user.GitHubID,
			"github_username":  user.GitHubUsername,
			"google_id":        user.GoogleID,
			"google_email":     user.GoogleEmail,
			"facebook_id":      user.FacebookID,
			"facebook_email":   user.FacebookEmail,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
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

func generateAccessToken(userID uint, email string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
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

func generateRefreshToken(userID uint) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	refreshToken := models.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := database.DB.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

func revokeRefreshToken(token string) error {
	return database.DB.Model(&models.RefreshToken{}).Where("token = ?", token).Update("is_revoked", true).Error
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
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"created_at":       user.CreatedAt,
			"updated_at":       user.UpdatedAt,
			"phone":            user.Phone,
			"birthday":         user.Birthday,
			"gender":           user.Gender,
			"country":          user.Country,
			"lang":             user.Lang,
			"login_provider":   user.LoginProvider,
			"profile_image":    user.ProfileImage,
			"github_id":        user.GitHubID,
			"github_username":  user.GitHubUsername,
			"google_id":        user.GoogleID,
			"google_email":     user.GoogleEmail,
			"facebook_id":      user.FacebookID,
			"facebook_email":   user.FacebookEmail,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
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
		"message":         "GitHub account linked successfully",
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

func LinkGoogleAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GoogleLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if googleClientID == "" || googleClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}

	redirectURI := strings.TrimSpace(req.RedirectURI)
	if redirectURI == "" {
		redirectURI = getRedirectURI("GOOGLE_REDIRECT_URI", "MOBILE_GOOGLE_REDIRECT_URI", "http://localhost:3000/callback", false)
	}

	tokenResp, err := exchangeGoogleCodeForToken(req.Code, googleClientID, googleClientSecret, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	googleUser, err := getGoogleUser(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Google user"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("google_id = ?", googleUser.ID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This Google account is already linked to another user"})
		return
	}

	user.GoogleID = &googleUser.ID
	user.GoogleEmail = &googleUser.Email

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Google account"})
		return
	}

	var oauth2Token models.OAuth2Token
	err = database.DB.Where("user_id = ? AND service = ?", userID, "google").First(&oauth2Token).Error

	expiry := time.Now()
	if tokenResp.ExpiresIn > 0 {
		expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	} else {
		expiry = time.Now().Add(time.Hour)
	}

	if err != nil {
		oauth2Token = models.OAuth2Token{
			UserID:       userID.(uint),
			Service:      "google",
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
			TokenType:    tokenResp.TokenType,
			ExpiresAt:    &expiry,
			Scope:        tokenResp.Scope,
		}

		if err := database.DB.Create(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store Google OAuth2 token"})
			return
		}
	} else {
		oauth2Token.AccessToken = tokenResp.AccessToken
		if tokenResp.RefreshToken != "" {
			oauth2Token.RefreshToken = tokenResp.RefreshToken
		}
		oauth2Token.TokenType = tokenResp.TokenType
		oauth2Token.ExpiresAt = &expiry
		oauth2Token.Scope = tokenResp.Scope

		if err := database.DB.Save(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Google OAuth2 token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Google account linked successfully",
		"google_email": user.GoogleEmail,
	})
}

func LinkSpotifyAccount(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	var req SpotifyLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_LINK_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "https://overeasily-superable-catarina.ngrok-free.dev/oauth2/spotify/callback"
	}

	if spotifyClientID == "" || spotifyClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Spotify OAuth not configured"})
		return
	}

	tokenResp, err := exchangeSpotifyCodeForToken(req.Code, spotifyClientID, spotifyClientSecret, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for Spotify token"})
		return
	}

	spotifyUser, err := getSpotifyUser(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Spotify user"})
		return
	}

	if spotifyUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spotify account does not provide an email address"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existing models.User
	if err := database.DB.Where("spotify_id = ? AND id <> ?", spotifyUser.ID, uid).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This Spotify account is already linked to another user"})
		return
	}

	user.SpotifyID = &spotifyUser.ID
	user.SpotifyEmail = &spotifyUser.Email
	if user.LoginProvider == "" || user.LoginProvider == "email" {
		user.LoginProvider = "spotify"
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Spotify account"})
		return
	}

	expiry := time.Now().Add(time.Hour)
	if tokenResp.ExpiresIn > 0 {
		expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	var oauth2Token models.OAuth2Token
	findErr := database.DB.Where("user_id = ? AND service = ?", uid, "spotify").First(&oauth2Token).Error
	if findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			oauth2Token = models.OAuth2Token{
				UserID:       uid,
				Service:      "spotify",
				AccessToken:  tokenResp.AccessToken,
				RefreshToken: tokenResp.RefreshToken,
				TokenType:    tokenResp.TokenType,
				Scope:        tokenResp.Scope,
				ExpiresAt:    &expiry,
			}

			if err := database.DB.Create(&oauth2Token).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store Spotify OAuth2 token"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Spotify OAuth2 token"})
			return
		}
	} else {
		oauth2Token.AccessToken = tokenResp.AccessToken
		if tokenResp.RefreshToken != "" {
			oauth2Token.RefreshToken = tokenResp.RefreshToken
		}
		oauth2Token.TokenType = tokenResp.TokenType
		oauth2Token.Scope = tokenResp.Scope
		oauth2Token.ExpiresAt = &expiry

		if err := database.DB.Save(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Spotify OAuth2 token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Spotify account linked successfully",
		"spotify_email": user.SpotifyEmail,
	})
}

func UnlinkSpotifyAccount(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.SpotifyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Spotify account linked"})
		return
	}

	user.SpotifyID = nil
	user.SpotifyEmail = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink Spotify account"})
		return
	}

	if err := database.DB.Where("user_id = ? AND service = ?", uid, "spotify").Delete(&models.OAuth2Token{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove Spotify tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Spotify account unlinked successfully",
	})
}

func UnlinkGoogleAccount(c *gin.Context) {
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

	if user.GoogleID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Google account linked"})
		return
	}

	user.GoogleID = nil
	user.GoogleEmail = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink Google account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Google account unlinked successfully",
	})
}

func exchangeGoogleCodeForToken(code, clientID, clientSecret, redirectURI string) (*GoogleTokenResponse, error) {
	url := "https://oauth2.googleapis.com/token"

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectURI,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

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

	var tokenResp GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getGoogleUser(accessToken string) (*GoogleUserResponse, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

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

	var googleUser GoogleUserResponse
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func exchangeSpotifyCodeForToken(code, clientID, clientSecret, redirectURI string) (*SpotifyTokenResponse, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", "Basic "+basicAuth)

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify token exchange failed: %s", string(body))
	}

	var tokenResp SpotifyTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getSpotifyUser(accessToken string) (*SpotifyUserProfile, error) {
	userURL := "https://api.spotify.com/v1/me"

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify user info request failed: %s", string(body))
	}

	var profile SpotifyUserProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

type GitHubRepository struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PushedAt      string `json:"pushed_at"`
}

type GitHubRepositoriesResponse struct {
	Repositories []GitHubRepository `json:"repositories"`
}

func GetGitHubRepositories(c *gin.Context) {
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

	if user.GitHubUsername == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No GitHub account linked"})
		return
	}

	fmt.Printf("🔍 Fetching repositories for user: %s (GitHub username: %s)\n", user.Email, *user.GitHubUsername)

	repositories, err := getGitHubRepositoriesForUser(*user.GitHubUsername)
	if err != nil {
		fmt.Printf("❌ Error fetching repositories: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repositories"})
		return
	}

	fmt.Printf("✅ Successfully fetched %d repositories\n", len(repositories))
	c.JSON(http.StatusOK, GitHubRepositoriesResponse{
		Repositories: repositories,
	})
}

func getGitHubRepositoriesForUser(username string) ([]GitHubRepository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?type=public&per_page=100", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "AREA-App")

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken != "" {
		req.Header.Set("Authorization", "token "+githubToken)
		fmt.Printf("🔑 Using GitHub token for API request\n")
	} else {
		fmt.Printf("⚠️ No GitHub token found, using public API (rate limited)\n")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var githubRepos []GitHubRepository
	if err := json.Unmarshal(body, &githubRepos); err != nil {
		return nil, err
	}

	return githubRepos, nil
}

type GitHubGmailAreaRequest struct {
	RepositoryID      int      `json:"repository_id" binding:"required"`
	DestinationEmail  string   `json:"destination_email" binding:"required,email"`
	NotificationTypes []string `json:"notification_types" binding:"required"`
}

func CreateGitHubGmailArea(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GitHubGmailAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.GitHubUsername == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GitHub username not configured. Please link your GitHub account first."})
		return
	}

	repositories, err := getGitHubRepositoriesForUser(*user.GitHubUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository information"})
		return
	}

	var targetRepo *GitHubRepository
	for _, repo := range repositories {
		if repo.ID == req.RepositoryID {
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository not found"})
		return
	}

	emailService, err := services.NewEmailService()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to initialize email service"})
		return
	}

	triggerConfig := map[string]interface{}{
		"repository_id":        req.RepositoryID,
		"notification_types":   req.NotificationTypes,
		"repository_name":      targetRepo.Name,
		"repository_full_name": targetRepo.FullName,
	}

	actionConfig := map[string]interface{}{
		"destination_email": req.DestinationEmail,
		"subject_template":  emailService.GetDefaultPushSubjectTemplate(),
		"body_template":     emailService.GetDefaultPushBodyTemplate(),
	}

	triggerConfigJSON, _ := json.Marshal(triggerConfig)
	actionConfigJSON, _ := json.Marshal(actionConfig)

	area := models.Area{
		UserID:         userID.(uint),
		Name:           fmt.Sprintf("GitHub → Gmail (%s)", targetRepo.Name),
		Description:    fmt.Sprintf("Envoie des emails Gmail lors d'événements sur le repository %s", targetRepo.FullName),
		IsActive:       true,
		TriggerService: "github",
		TriggerType:    "push",
		TriggerConfig:  datatypes.JSON(triggerConfigJSON),
		ActionService:  "gmail",
		ActionType:     "send_email",
		ActionConfig:   datatypes.JSON(actionConfigJSON),
	}

	if err := database.DB.Create(&area).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create area"})
		return
	}

	githubService := services.NewGitHubIntegrationService()
	webhookResp, err := githubService.CreateWebhook(targetRepo.FullName[:strings.Index(targetRepo.FullName, "/")], targetRepo.Name)

	var webhookMessage string
	if err != nil {
		webhookMessage = fmt.Sprintf("Area created but webhook configuration failed: %v", err)
	} else {
		webhookMessage = fmt.Sprintf("Webhook configured successfully (ID: %d)", webhookResp.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "GitHub-Gmail area created successfully",
		"area_id":            area.ID,
		"repository_id":      req.RepositoryID,
		"repository_name":    targetRepo.Name,
		"destination_email":  req.DestinationEmail,
		"notification_types": req.NotificationTypes,
		"webhook_status":     webhookMessage,
	})
}

func LinkFacebookAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req FacebookLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	facebookClientID := os.Getenv("FACEBOOK_CLIENT_ID")
	facebookClientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")

	if facebookClientID == "" || facebookClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Facebook OAuth not configured"})
		return
	}

	redirectURI := getRedirectURI("FACEBOOK_REDIRECT_URI", "MOBILE_FACEBOOK_REDIRECT_URI", "http://localhost:3000/auth/facebook/callback", false)
	accessToken, err := exchangeFacebookCodeForToken(req.Code, facebookClientID, facebookClientSecret, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	facebookUser, err := getFacebookUser(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Facebook user"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("facebook_id = ?", facebookUser.ID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This Facebook account is already linked to another user"})
		return
	}

	user.FacebookID = &facebookUser.ID
	user.FacebookEmail = &facebookUser.Email

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Facebook account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Facebook account linked successfully",
		"facebook_email": user.FacebookEmail,
	})
}

func UnlinkFacebookAccount(c *gin.Context) {
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

	if user.FacebookID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Facebook account linked"})
		return
	}

	user.FacebookID = nil
	user.FacebookEmail = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink Facebook account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Facebook account unlinked successfully",
	})
}

func exchangeFacebookCodeForToken(code, clientID, clientSecret, redirectURI string) (string, error) {
	url := "https://graph.facebook.com/v18.0/oauth/access_token"

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
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

	var tokenResp FacebookTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func getFacebookUser(accessToken string) (*FacebookUserResponse, error) {
	url := "https://graph.facebook.com/v18.0/me?fields=id,name"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

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

	var facebookUser FacebookUserResponse
	if err := json.Unmarshal(body, &facebookUser); err != nil {
		return nil, err
	}

	if facebookUser.Email == "" {
		facebookUser.Email = facebookUser.ID + "@facebook.com"
	}

	return &facebookUser, nil
}

func MobileOAuth2Login(c *gin.Context) {
	var req OAuth2LoginRequest

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

	accessToken, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	c.JSON(http.StatusOK, OAuth2TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: gin.H{
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
		},
	})
}

func GitHubDirectLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code required"})
		return
	}

	state := c.Query("state")
	stateInfo := parseMobileState(state)

	if stateInfo.isMobile && stateInfo.mode == "link" {
		redirectToMobileWithAuthorizationCode(c, "github", code, "")
		return
	}

	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	if githubClientID == "" || githubClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub OAuth not configured"})
		return
	}

	accessToken, err := exchangeCodeForToken(code, githubClientID, githubClientSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	githubUser, err := getGitHubUser(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get GitHub user"})
		return
	}

	githubIDStr := fmt.Sprintf("%d", githubUser.ID)

	var user models.User
	if err := database.DB.Where("github_id = ?", githubIDStr).First(&user).Error; err != nil {
		user = models.User{
			Email:          githubUser.Email,
			FirstName:      githubUser.Name,
			GitHubID:       &githubIDStr,
			GitHubUsername: &githubUser.Login,
			LoginProvider:  "github",
			IsActive:       true,
			Role:           "member",
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	accessTokenJWT, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	tokenResponse := OAuth2TokenResponse{
		AccessToken:  accessTokenJWT,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"profile_image":   user.ProfileImage,
			"role":            user.Role,
			"is_active":       user.IsActive,
			"github_id":       user.GitHubID,
			"github_username": user.GitHubUsername,
			"spotify_id":      user.SpotifyID,
			"spotify_email":   user.SpotifyEmail,
			"google_id":       user.GoogleID,
			"google_email":    user.GoogleEmail,
			"facebook_id":     user.FacebookID,
			"facebook_email":  user.FacebookEmail,
		},
	}

	if stateInfo.isMobile {
		redirectToMobileWithTokens(c, "github", tokenResponse)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func GoogleDirectLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code required"})
		return
	}

	state := c.Query("state")
	stateInfo := parseMobileState(state)

	if stateInfo.isMobile && stateInfo.mode == "link" {
		redirectToMobileWithAuthorizationCode(c, "google", code, "")
		return
	}

	redirectURI := getRedirectURI("GOOGLE_REDIRECT_URI", "MOBILE_GOOGLE_REDIRECT_URI", "http://localhost:3000/callback", stateInfo.isMobile)

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if googleClientID == "" || googleClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}

	tokenResp, err := exchangeGoogleCodeForToken(code, googleClientID, googleClientSecret, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	googleUser, err := getGoogleUser(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Google user"})
		return
	}

	var user models.User
	if err := database.DB.Where("google_id = ?", googleUser.ID).First(&user).Error; err != nil {
		user = models.User{
			Email:         googleUser.Email,
			FirstName:     googleUser.GivenName,
			LastName:      googleUser.FamilyName,
			GoogleID:      &googleUser.ID,
			GoogleEmail:   &googleUser.Email,
			LoginProvider: "google",
			IsActive:      true,
			Role:          "member",
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	accessTokenJWT, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	var oauth2Token models.OAuth2Token
	err = database.DB.Where("user_id = ? AND service = ?", user.ID, "google").First(&oauth2Token).Error

	expiry := time.Now()
	if tokenResp.ExpiresIn > 0 {
		expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	} else {
		expiry = time.Now().Add(time.Hour)
	}

	if err != nil {
		oauth2Token = models.OAuth2Token{
			UserID:       user.ID,
			Service:      "google",
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
			TokenType:    tokenResp.TokenType,
			ExpiresAt:    &expiry,
			Scope:        tokenResp.Scope,
		}

		if err := database.DB.Create(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store Google OAuth2 token"})
			return
		}
	} else {
		oauth2Token.AccessToken = tokenResp.AccessToken
		if tokenResp.RefreshToken != "" {
			oauth2Token.RefreshToken = tokenResp.RefreshToken
		}
		oauth2Token.TokenType = tokenResp.TokenType
		oauth2Token.ExpiresAt = &expiry
		oauth2Token.Scope = tokenResp.Scope

		if err := database.DB.Save(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Google OAuth2 token"})
			return
		}
	}

	tokenResponse := OAuth2TokenResponse{
		AccessToken:  accessTokenJWT,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"profile_image":   user.ProfileImage,
			"role":            user.Role,
			"is_active":       user.IsActive,
			"google_id":       user.GoogleID,
			"google_email":    user.GoogleEmail,
			"spotify_id":      user.SpotifyID,
			"spotify_email":   user.SpotifyEmail,
			"github_id":       user.GitHubID,
			"github_username": user.GitHubUsername,
			"facebook_id":     user.FacebookID,
			"facebook_email":  user.FacebookEmail,
		},
	}

	if stateInfo.isMobile {
		redirectToMobileWithTokens(c, "google", tokenResponse)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func SpotifyDirectLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code required"})
		return
	}

	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	state := c.Query("state")
	stateInfo := parseMobileState(state)

	if stateInfo.isMobile && stateInfo.mode == "link" {
		redirectToMobileWithAuthorizationCode(c, "spotify", code, "")
		return
	}

	redirectURI := getRedirectURI("SPOTIFY_REDIRECT_URI", "MOBILE_SPOTIFY_REDIRECT_URI", "http://127.0.0.1:3000/oauth2/spotify/callback", stateInfo.isMobile)

	if spotifyClientID == "" || spotifyClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Spotify OAuth not configured"})
		return
	}

	tokenResp, err := exchangeSpotifyCodeForToken(code, spotifyClientID, spotifyClientSecret, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for Spotify token"})
		return
	}

	spotifyUser, err := getSpotifyUser(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Spotify user"})
		return
	}

	if spotifyUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spotify account does not provide an email address"})
		return
	}

	var user models.User
	result := database.DB.Where("spotify_id = ?", spotifyUser.ID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := database.DB.Where("email = ?", spotifyUser.Email).First(&user).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to look up user"})
					return
				}

				user = models.User{
					Email:         spotifyUser.Email,
					FirstName:     spotifyUser.DisplayName,
					SpotifyID:     &spotifyUser.ID,
					SpotifyEmail:  &spotifyUser.Email,
					LoginProvider: "spotify",
					IsActive:      true,
					Role:          "member",
				}

				if err := database.DB.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}
	}

	if user.SpotifyID == nil || *user.SpotifyID != spotifyUser.ID {
		user.SpotifyID = &spotifyUser.ID
	}
	user.SpotifyEmail = &spotifyUser.Email
	if user.LoginProvider == "" || user.LoginProvider == "email" {
		user.LoginProvider = "spotify"
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	expiry := time.Now().Add(time.Hour)
	if tokenResp.ExpiresIn > 0 {
		expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	var oauth2Token models.OAuth2Token
	tokenResult := database.DB.Where("user_id = ? AND service = ?", user.ID, "spotify").First(&oauth2Token)
	if tokenResult.Error != nil {
		if errors.Is(tokenResult.Error, gorm.ErrRecordNotFound) {
			oauth2Token = models.OAuth2Token{
				UserID:       user.ID,
				Service:      "spotify",
				AccessToken:  tokenResp.AccessToken,
				RefreshToken: tokenResp.RefreshToken,
				TokenType:    tokenResp.TokenType,
				Scope:        tokenResp.Scope,
				ExpiresAt:    &expiry,
			}

			if err := database.DB.Create(&oauth2Token).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store Spotify OAuth2 token"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Spotify OAuth2 token"})
			return
		}
	} else {
		oauth2Token.AccessToken = tokenResp.AccessToken
		if tokenResp.RefreshToken != "" {
			oauth2Token.RefreshToken = tokenResp.RefreshToken
		}
		oauth2Token.TokenType = tokenResp.TokenType
		oauth2Token.Scope = tokenResp.Scope
		oauth2Token.ExpiresAt = &expiry

		if err := database.DB.Save(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Spotify OAuth2 token"})
			return
		}
	}

	accessTokenJWT, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	tokenResponse := OAuth2TokenResponse{
		AccessToken:  accessTokenJWT,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"profile_image":   user.ProfileImage,
			"role":            user.Role,
			"is_active":       user.IsActive,
			"spotify_id":      user.SpotifyID,
			"spotify_email":   user.SpotifyEmail,
			"github_id":       user.GitHubID,
			"github_username": user.GitHubUsername,
			"google_id":       user.GoogleID,
			"google_email":    user.GoogleEmail,
			"facebook_id":     user.FacebookID,
			"facebook_email":  user.FacebookEmail,
		},
	}

	if stateInfo.isMobile {
		redirectToMobileWithTokens(c, "spotify", tokenResponse)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func FacebookDirectLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code required"})
		return
	}

	state := c.Query("state")
	stateInfo := parseMobileState(state)

	if stateInfo.isMobile && stateInfo.mode == "link" {
		redirectToMobileWithAuthorizationCode(c, "facebook", code, "")
		return
	}

	redirectURI := getRedirectURI("FACEBOOK_REDIRECT_URI", "MOBILE_FACEBOOK_REDIRECT_URI", "http://localhost:3000/auth/facebook/callback", stateInfo.isMobile)

	facebookClientID := os.Getenv("FACEBOOK_CLIENT_ID")
	facebookClientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")

	if facebookClientID == "" || facebookClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Facebook OAuth not configured"})
		return
	}

	accessToken, err := exchangeFacebookCodeForToken(code, facebookClientID, facebookClientSecret, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	facebookUser, err := getFacebookUser(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Facebook user"})
		return
	}

	var user models.User
	if err := database.DB.Where("facebook_id = ?", facebookUser.ID).First(&user).Error; err != nil {
		user = models.User{
			Email:         facebookUser.Email,
			FirstName:     facebookUser.Name,
			FacebookID:    &facebookUser.ID,
			FacebookEmail: &facebookUser.Email,
			LoginProvider: "facebook",
			IsActive:      true,
			Role:          "member",
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	accessTokenJWT, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	tokenResponse := OAuth2TokenResponse{
		AccessToken:  accessTokenJWT,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: gin.H{
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"profile_image":    user.ProfileImage,
			"role":             user.Role,
			"is_active":        user.IsActive,
			"facebook_id":      user.FacebookID,
			"facebook_email":   user.FacebookEmail,
			"spotify_id":       user.SpotifyID,
			"spotify_email":    user.SpotifyEmail,
			"github_id":        user.GitHubID,
			"github_username":  user.GitHubUsername,
			"google_id":        user.GoogleID,
			"google_email":     user.GoogleEmail,
			"twitter_id":       user.TwitterID,
			"twitter_username": user.TwitterUsername,
		},
	}

	if stateInfo.isMobile {
		redirectToMobileWithTokens(c, "facebook", tokenResponse)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// Twitter OAuth 2.0 helpers
func exchangeTwitterCodeForToken(code, clientID, clientSecret, redirectURI, codeVerifier string) (*TwitterTokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", "https://api.twitter.com/2/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitter token exchange failed: %s", string(body))
	}

	var tokenResp TwitterTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getTwitterUser(accessToken string) (*TwitterUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get twitter user: %s", string(body))
	}

	var userResp TwitterUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}

// LinkTwitterAccount links a Twitter account to the logged-in user
func LinkTwitterAccount(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	var req TwitterLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	twitterClientID := os.Getenv("TWITTER_CLIENT_ID")
	twitterClientSecret := os.Getenv("TWITTER_CLIENT_SECRET")
	redirectURI := os.Getenv("TWITTER_REDIRECT_URI")

	if redirectURI == "" {
		redirectURI = "https://overeasily-superable-catarina.ngrok-free.dev/oauth2/twitter/callback"
	}

	if twitterClientID == "" || twitterClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Twitter OAuth not configured"})
		return
	}

	tokenResp, err := exchangeTwitterCodeForToken(req.Code, twitterClientID, twitterClientSecret, redirectURI, req.CodeVerifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for Twitter token"})
		return
	}

	twitterUser, err := getTwitterUser(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Twitter user"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existing models.User
	if err := database.DB.Where("twitter_id = ? AND id <> ?", twitterUser.Data.ID, uid).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This Twitter account is already linked to another user"})
		return
	}

	user.TwitterID = &twitterUser.Data.ID
	user.TwitterUsername = &twitterUser.Data.Username

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Twitter account"})
		return
	}

	expiry := time.Now().Add(time.Hour * 2)
	if tokenResp.ExpiresIn > 0 {
		expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	var oauth2Token models.OAuth2Token
	findErr := database.DB.Where("user_id = ? AND service = ?", uid, "twitter").First(&oauth2Token).Error
	if findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			oauth2Token = models.OAuth2Token{
				UserID:       uid,
				Service:      "twitter",
				AccessToken:  tokenResp.AccessToken,
				RefreshToken: tokenResp.RefreshToken,
				TokenType:    tokenResp.TokenType,
				Scope:        tokenResp.Scope,
				ExpiresAt:    &expiry,
			}

			if err := database.DB.Create(&oauth2Token).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store Twitter OAuth2 token"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Twitter OAuth2 token"})
			return
		}
	} else {
		oauth2Token.AccessToken = tokenResp.AccessToken
		if tokenResp.RefreshToken != "" {
			oauth2Token.RefreshToken = tokenResp.RefreshToken
		}
		oauth2Token.TokenType = tokenResp.TokenType
		oauth2Token.Scope = tokenResp.Scope
		oauth2Token.ExpiresAt = &expiry

		if err := database.DB.Save(&oauth2Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Twitter OAuth2 token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Twitter account linked successfully",
		"twitter_username": user.TwitterUsername,
	})
}

// UnlinkTwitterAccount unlinks Twitter from the current user
func UnlinkTwitterAccount(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.TwitterID = nil
	user.TwitterUsername = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink Twitter account"})
		return
	}

	database.DB.Where("user_id = ? AND service = ?", uid, "twitter").Delete(&models.OAuth2Token{})

	c.JSON(http.StatusOK, gin.H{"message": "Twitter account unlinked successfully"})
}

// TwitterDirectLogin handles Twitter OAuth callback for direct login
func TwitterDirectLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, "https://overeasily-superable-catarina.ngrok-free.dev/login?error=missing_code")
		return
	}

	state := c.Query("state")
	if state == "link" {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://overeasily-superable-catarina.ngrok-free.dev/auth/twitter/callback?code=%s&state=link", code))
		return
	}
	stateInfo := parseMobileState(state)
	if stateInfo.isMobile && stateInfo.mode == "link" {
		redirectToMobileWithAuthorizationCode(c, "twitter", code, stateInfo.payload)
		return
	}

	codeVerifier := ""
	if stateInfo.isMobile && stateInfo.payload != "" {
		codeVerifier = stateInfo.payload
	}

	twitterClientID := os.Getenv("TWITTER_CLIENT_ID")
	twitterClientSecret := os.Getenv("TWITTER_CLIENT_SECRET")
	redirectURI := getRedirectURI("TWITTER_REDIRECT_URI", "MOBILE_TWITTER_REDIRECT_URI", "https://overeasily-superable-catarina.ngrok-free.dev/oauth2/twitter/callback", stateInfo.isMobile)

	if twitterClientID == "" || twitterClientSecret == "" {
		c.Redirect(http.StatusTemporaryRedirect, "https://overeasily-superable-catarina.ngrok-free.dev/login?error=twitter_not_configured")
		return
	}

	// Note: Direct login doesn't support PKCE in this implementation
	// For production, implement PKCE for direct login as well
	tokenResp, err := exchangeTwitterCodeForToken(code, twitterClientID, twitterClientSecret, redirectURI, codeVerifier)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://overeasily-superable-catarina.ngrok-free.dev/login?error=token_exchange_failed")
		return
	}

	twitterUser, err := getTwitterUser(tokenResp.AccessToken)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://overeasily-superable-catarina.ngrok-free.dev/login?error=user_fetch_failed")
		return
	}

	var user models.User
	if err := database.DB.Where("twitter_id = ?", twitterUser.Data.ID).First(&user).Error; err != nil {
		user = models.User{
			Email:           twitterUser.Data.Username + "@twitter.placeholder",
			FirstName:       twitterUser.Data.Name,
			TwitterID:       &twitterUser.Data.ID,
			TwitterUsername: &twitterUser.Data.Username,
			LoginProvider:   "twitter",
			IsActive:        true,
			Role:            "member",
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "http://127.0.0.1:3000/login?error=user_creation_failed")
			return
		}
	}

	accessTokenJWT, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "http://127.0.0.1:3000/login?error=token_generation_failed")
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "http://127.0.0.1:3000/login?error=refresh_token_failed")
		return
	}

	if stateInfo.isMobile {
		tokenResponse := OAuth2TokenResponse{
			AccessToken:  accessTokenJWT,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    900,
			User: gin.H{
				"id":               user.ID,
				"email":            user.Email,
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"profile_image":    user.ProfileImage,
				"role":             user.Role,
				"is_active":        user.IsActive,
				"twitter_id":       user.TwitterID,
				"twitter_username": user.TwitterUsername,
				"spotify_id":       user.SpotifyID,
				"spotify_email":    user.SpotifyEmail,
				"github_id":        user.GitHubID,
				"github_username":  user.GitHubUsername,
				"google_id":        user.GoogleID,
				"google_email":     user.GoogleEmail,
				"facebook_id":      user.FacebookID,
				"facebook_email":   user.FacebookEmail,
			},
		}

		redirectToMobileWithTokens(c, "twitter", tokenResponse)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(
		"http://127.0.0.1:3000/home?access_token=%s&refresh_token=%s&token_type=Bearer",
		accessTokenJWT,
		refreshToken,
	))
}
