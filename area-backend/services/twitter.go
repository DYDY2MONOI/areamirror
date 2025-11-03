package services

import (
	"Golang-API-tutoriel/database"
	"Golang-API-tutoriel/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	twitterAPIBaseURL    = "https://api.twitter.com/2"
	twitterOAuthTokenURL = "https://api.twitter.com/2/oauth2/token"
	defaultTwitterUA     = "AREA-Automation/1.0"
)

type TwitterService struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

type TwitterTweet struct {
	ID              string
	Text            string
	AuthorID        string
	AuthorUsername  string
	AuthorName      string
	CreatedAt       time.Time
	URL             string
	ConversationID  string
	InReplyToUserID string
}

type twitterMentionResponse struct {
	Data []struct {
		ID              string `json:"id"`
		Text            string `json:"text"`
		AuthorID        string `json:"author_id"`
		CreatedAt       string `json:"created_at"`
		ConversationID  string `json:"conversation_id"`
		InReplyToUserID string `json:"in_reply_to_user_id"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
	Meta struct {
		NewestID    string `json:"newest_id"`
		OldestID    string `json:"oldest_id"`
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
	} `json:"meta"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type twitterUserTweetsResponse struct {
	Data []struct {
		ID            string `json:"id"`
		Text          string `json:"text"`
		CreatedAt     string `json:"created_at"`
		PublicMetrics struct {
			LikeCount    int `json:"like_count"`
			RetweetCount int `json:"retweet_count"`
			ReplyCount   int `json:"reply_count"`
			QuoteCount   int `json:"quote_count"`
		} `json:"public_metrics"`
	} `json:"data"`
	Meta struct {
		NewestID    string `json:"newest_id"`
		OldestID    string `json:"oldest_id"`
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
	} `json:"meta"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type twitterFollowersResponse struct {
	Data []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Username    string `json:"username"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
	} `json:"data"`
	Meta struct {
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
	} `json:"meta"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type twitterTweetCreateResponse struct {
	Data struct {
		ID             string `json:"id"`
		Text           string `json:"text"`
		ConversationID string `json:"conversation_id"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type twitterRetweetResponse struct {
	Data struct {
		Retweeted bool `json:"retweeted"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type TwitterTweetMetrics struct {
	ID           string
	Text         string
	CreatedAt    time.Time
	LikeCount    int
	RetweetCount int
	ReplyCount   int
	QuoteCount   int
}

type TwitterUser struct {
	ID          string
	Username    string
	Name        string
	Description string
	CreatedAt   time.Time
}

type twitterTokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type TwitterAPIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *TwitterAPIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("twitter api error (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("twitter api error (%d)", e.StatusCode)
}

func NewTwitterService() (*TwitterService, error) {
	clientID := strings.TrimSpace(os.Getenv("TWITTER_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("TWITTER_CLIENT_SECRET"))

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("Twitter OAuth credentials not configured")
	}

	return &TwitterService{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

func (s *TwitterService) getTwitterToken(userID uint) (*models.OAuth2Token, error) {
	var token models.OAuth2Token
	if err := database.DB.Where("user_id = ? AND service = ?", userID, "twitter").First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no Twitter account linked for user %d", userID)
		}
		return nil, fmt.Errorf("failed to fetch twitter token for user %d: %w", userID, err)
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("twitter access token missing for user %d", userID)
	}

	if token.NeedsRefresh() {
		if err := s.refreshTwitterToken(&token); err != nil {
			return nil, fmt.Errorf("failed to refresh twitter token for user %d: %w", userID, err)
		}
	}

	return &token, nil
}

func (s *TwitterService) refreshTwitterToken(token *models.OAuth2Token) error {
	if token.RefreshToken == "" {
		return fmt.Errorf("no refresh token stored for twitter user %d", token.UserID)
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", token.RefreshToken)

	req, err := http.NewRequest(http.MethodPost, twitterOAuthTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", defaultTwitterUA)
	req.SetBasicAuth(s.clientID, s.clientSecret)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh twitter token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read twitter token refresh response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    extractTwitterErrorMessage(body),
			Body:       string(body),
		}
	}

	var refreshResp twitterTokenRefreshResponse
	if err := json.Unmarshal(body, &refreshResp); err != nil {
		return fmt.Errorf("failed to parse twitter token refresh response: %w", err)
	}

	if refreshResp.AccessToken == "" {
		return fmt.Errorf("twitter refresh response missing access token")
	}

	token.AccessToken = refreshResp.AccessToken
	token.TokenType = refreshResp.TokenType
	token.Scope = refreshResp.Scope

	if refreshResp.RefreshToken != "" {
		token.RefreshToken = refreshResp.RefreshToken
	}

	if refreshResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(refreshResp.ExpiresIn) * time.Second)
		token.ExpiresAt = &expiry
	}

	if err := database.DB.Save(token).Error; err != nil {
		return fmt.Errorf("failed to persist refreshed twitter token: %w", err)
	}

	return nil
}

func (s *TwitterService) FetchMentions(userID uint, twitterUserID, sinceID string) ([]TwitterTweet, string, error) {
	token, err := s.getTwitterToken(userID)
	if err != nil {
		return nil, "", err
	}

	tweets, newestID, err := s.fetchMentionsWithToken(token, twitterUserID, sinceID)
	if err != nil {
		var apiErr *TwitterAPIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			if err := s.refreshTwitterToken(token); err != nil {
				return nil, "", fmt.Errorf("failed to refresh twitter token after unauthorized response: %w", err)
			}
			return s.fetchMentionsWithToken(token, twitterUserID, sinceID)
		}
		return nil, "", err
	}

	return tweets, newestID, nil
}

func (s *TwitterService) fetchMentionsWithToken(token *models.OAuth2Token, twitterUserID, sinceID string) ([]TwitterTweet, string, error) {
	if strings.TrimSpace(twitterUserID) == "" {
		return nil, "", fmt.Errorf("twitter user id is required to fetch mentions")
	}

	endpoint := fmt.Sprintf("%s/users/%s/mentions", twitterAPIBaseURL, twitterUserID)
	params := url.Values{}
	params.Set("max_results", "5")
	params.Set("tweet.fields", "created_at,author_id,conversation_id,in_reply_to_user_id")
	params.Set("expansions", "author_id")
	params.Set("user.fields", "username,name")
	if sinceID != "" {
		params.Set("since_id", sinceID)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", defaultTwitterUA)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode >= 400 {
		return nil, "", &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    extractTwitterErrorMessage(body),
			Body:       string(body),
		}
	}

	var apiResp twitterMentionResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode twitter mentions response: %w", err)
	}

	if apiResp.Meta.ResultCount == 0 || len(apiResp.Data) == 0 {
		return []TwitterTweet{}, sinceID, nil
	}

	authors := map[string]struct {
		username string
		name     string
	}{}
	for _, user := range apiResp.Includes.Users {
		authors[user.ID] = struct {
			username string
			name     string
		}{
			username: user.Username,
			name:     user.Name,
		}
	}

	tweets := make([]TwitterTweet, 0, len(apiResp.Data))
	for _, item := range apiResp.Data {
		createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
		if err != nil {
			createdAt = time.Now().UTC()
		}

		userInfo := authors[item.AuthorID]
		tweet := TwitterTweet{
			ID:              item.ID,
			Text:            item.Text,
			AuthorID:        item.AuthorID,
			AuthorUsername:  userInfo.username,
			AuthorName:      userInfo.name,
			CreatedAt:       createdAt,
			ConversationID:  item.ConversationID,
			InReplyToUserID: item.InReplyToUserID,
		}
		if tweet.AuthorUsername != "" && tweet.ID != "" {
			tweet.URL = fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.AuthorUsername, tweet.ID)
		}
		tweets = append(tweets, tweet)
	}

	for i, j := 0, len(tweets)-1; i < j; i, j = i+1, j-1 {
		tweets[i], tweets[j] = tweets[j], tweets[i]
	}

	newestID := apiResp.Meta.NewestID
	if newestID == "" && len(tweets) > 0 {
		newestID = tweets[len(tweets)-1].ID
	}

	return tweets, newestID, nil
}

func (s *TwitterService) FetchUserTweetsWithMetrics(userID uint, twitterUserID string, maxResults int) ([]TwitterTweetMetrics, error) {
	token, err := s.getTwitterToken(userID)
	if err != nil {
		return nil, err
	}

	tweets, err := s.fetchUserTweetsWithToken(token, twitterUserID, maxResults)
	if err != nil {
		var apiErr *TwitterAPIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			if err := s.refreshTwitterToken(token); err != nil {
				return nil, fmt.Errorf("failed to refresh twitter token after unauthorized response: %w", err)
			}
			return s.fetchUserTweetsWithToken(token, twitterUserID, maxResults)
		}
		return nil, err
	}

	return tweets, nil
}

func (s *TwitterService) fetchUserTweetsWithToken(token *models.OAuth2Token, twitterUserID string, maxResults int) ([]TwitterTweetMetrics, error) {
	if strings.TrimSpace(twitterUserID) == "" {
		return nil, fmt.Errorf("twitter user id is required to fetch tweets")
	}

	if maxResults <= 0 {
		maxResults = 20
	}
	if maxResults > 100 {
		maxResults = 100
	}

	endpoint := fmt.Sprintf("%s/users/%s/tweets", twitterAPIBaseURL, twitterUserID)
	params := url.Values{}
	params.Set("max_results", strconv.Itoa(maxResults))
	params.Set("tweet.fields", "created_at,public_metrics")
	params.Set("exclude", "retweets")

	req, err := http.NewRequest(http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", defaultTwitterUA)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    extractTwitterErrorMessage(body),
			Body:       string(body),
		}
	}

	var apiResp twitterUserTweetsResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse twitter tweets response: %w", err)
	}

	if apiResp.Meta.ResultCount == 0 || len(apiResp.Data) == 0 {
		return []TwitterTweetMetrics{}, nil
	}

	tweets := make([]TwitterTweetMetrics, 0, len(apiResp.Data))
	for _, item := range apiResp.Data {
		createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
		if err != nil {
			createdAt = time.Now().UTC()
		}
		tweets = append(tweets, TwitterTweetMetrics{
			ID:           item.ID,
			Text:         item.Text,
			CreatedAt:    createdAt,
			LikeCount:    item.PublicMetrics.LikeCount,
			RetweetCount: item.PublicMetrics.RetweetCount,
			ReplyCount:   item.PublicMetrics.ReplyCount,
			QuoteCount:   item.PublicMetrics.QuoteCount,
		})
	}

	return tweets, nil
}

func (s *TwitterService) FetchFollowers(userID uint, twitterUserID string, maxResults int) ([]TwitterUser, error) {
	token, err := s.getTwitterToken(userID)
	if err != nil {
		return nil, err
	}

	users, err := s.fetchFollowersWithToken(token, twitterUserID, maxResults)
	if err != nil {
		var apiErr *TwitterAPIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			if err := s.refreshTwitterToken(token); err != nil {
				return nil, fmt.Errorf("failed to refresh twitter token after unauthorized response: %w", err)
			}
			return s.fetchFollowersWithToken(token, twitterUserID, maxResults)
		}
		return nil, err
	}

	return users, nil
}

func (s *TwitterService) fetchFollowersWithToken(token *models.OAuth2Token, twitterUserID string, maxResults int) ([]TwitterUser, error) {
	if strings.TrimSpace(twitterUserID) == "" {
		return nil, fmt.Errorf("twitter user id is required to fetch followers")
	}

	if maxResults <= 0 {
		maxResults = 20
	}
	if maxResults > 1000 {
		maxResults = 1000
	}

	endpoint := fmt.Sprintf("%s/users/%s/followers", twitterAPIBaseURL, twitterUserID)
	params := url.Values{}
	params.Set("max_results", strconv.Itoa(maxResults))
	params.Set("user.fields", "created_at,description")

	req, err := http.NewRequest(http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", defaultTwitterUA)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    extractTwitterErrorMessage(body),
			Body:       string(body),
		}
	}

	var apiResp twitterFollowersResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse twitter followers response: %w", err)
	}

	if apiResp.Meta.ResultCount == 0 || len(apiResp.Data) == 0 {
		return []TwitterUser{}, nil
	}

	followers := make([]TwitterUser, 0, len(apiResp.Data))
	for _, item := range apiResp.Data {
		createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
		if err != nil {
			createdAt = time.Now().UTC()
		}
		followers = append(followers, TwitterUser{
			ID:          item.ID,
			Username:    item.Username,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   createdAt,
		})
	}

	return followers, nil
}

func (s *TwitterService) PostTweet(userID uint, text string, inReplyToID string) (*TwitterTweet, error) {
	token, err := s.getTwitterToken(userID)
	if err != nil {
		return nil, err
	}

	tweet, err := s.postTweetWithToken(token, text, inReplyToID)
	if err != nil {
		var apiErr *TwitterAPIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			if err := s.refreshTwitterToken(token); err != nil {
				return nil, fmt.Errorf("failed to refresh twitter token after unauthorized response: %w", err)
			}
			return s.postTweetWithToken(token, text, inReplyToID)
		}
		return nil, err
	}

	return tweet, nil
}

func (s *TwitterService) postTweetWithToken(token *models.OAuth2Token, text string, inReplyToID string) (*TwitterTweet, error) {
	payload := map[string]interface{}{
		"text": text,
	}

	if inReplyToID != "" {
		payload["reply"] = map[string]string{
			"in_reply_to_tweet_id": inReplyToID,
		}
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, twitterAPIBaseURL+"/tweets", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", defaultTwitterUA)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    extractTwitterErrorMessage(body),
			Body:       string(body),
		}
	}

	var apiResp twitterTweetCreateResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse twitter post response: %w", err)
	}

	if len(apiResp.Errors) > 0 {
		return nil, &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Errors[0].Message,
			Body:       string(body),
		}
	}

	if apiResp.Data.ID == "" {
		return nil, fmt.Errorf("twitter post response missing tweet id")
	}

	result := &TwitterTweet{
		ID:             apiResp.Data.ID,
		Text:           apiResp.Data.Text,
		CreatedAt:      time.Now().UTC(),
		ConversationID: apiResp.Data.ConversationID,
	}

	return result, nil
}

func (s *TwitterService) Retweet(userID uint, twitterUserID string, tweetID string) error {
	twitterUserID = strings.TrimSpace(twitterUserID)
	if twitterUserID == "" {
		return fmt.Errorf("twitter user id is required for retweets")
	}

	token, err := s.getTwitterToken(userID)
	if err != nil {
		return err
	}

	if err := s.retweetWithToken(token, twitterUserID, tweetID); err != nil {
		var apiErr *TwitterAPIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			if refreshErr := s.refreshTwitterToken(token); refreshErr != nil {
				return fmt.Errorf("failed to refresh twitter token after unauthorized retweet: %w", refreshErr)
			}
			return s.retweetWithToken(token, twitterUserID, tweetID)
		}
		return err
	}

	return nil
}

func (s *TwitterService) retweetWithToken(token *models.OAuth2Token, twitterUserID string, tweetID string) error {
	tweetID = strings.TrimSpace(tweetID)
	if tweetID == "" {
		return fmt.Errorf("tweet id is required for retweets")
	}

	endpoint := fmt.Sprintf("%s/users/%s/retweets", twitterAPIBaseURL, twitterUserID)
	payload := map[string]string{
		"tweet_id": tweetID,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", defaultTwitterUA)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    extractTwitterErrorMessage(body),
			Body:       string(body),
		}
	}

	var apiResp twitterRetweetResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse twitter retweet response: %w", err)
	}

	if len(apiResp.Errors) > 0 {
		return &TwitterAPIError{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Errors[0].Message,
			Body:       string(body),
		}
	}

	return nil
}

func extractTwitterErrorMessage(body []byte) string {
	var payload struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		if payload.Detail != "" {
			return payload.Detail
		}
		if payload.Title != "" {
			return payload.Title
		}
		if len(payload.Errors) > 0 && payload.Errors[0].Message != "" {
			return payload.Errors[0].Message
		}
	}
	return ""
}
