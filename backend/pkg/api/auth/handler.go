package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sharithg/civet/internal/config"
	"github.com/sharithg/civet/internal/genai"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/internal/storage"
)

type authRepository struct {
	DB      *pgxpool.Pool
	Ctx     *context.Context
	Storage *storage.Storage
	Genai   genai.OpenAi
	Config  *config.Config
	Repo    *repository.Queries
}

func New(db *pgxpool.Pool, repo *repository.Queries, storage *storage.Storage, genai genai.OpenAi, config *config.Config, ctx *context.Context) *authRepository {
	return &authRepository{
		DB:      db,
		Ctx:     ctx,
		Storage: storage,
		Genai:   genai,
		Config:  config,
		Repo:    repo,
	}
}

func (a *authRepository) RefreshTokenHandler(c *gin.Context) {
	platform := c.Query("platform")
	if platform == "" {
		platform = "native"
	}

	var refreshToken string

	// 1. Native client: Form body
	if platform == "native" {
		refreshToken = c.PostForm("refresh_token")
	}

	// 2. Web: Cookie
	if platform == "web" && refreshToken == "" {
		if cookie, err := c.Cookie(a.Config.RefreshCookieName); err == nil {
			refreshToken = cookie
		}
	}

	// 3. Fallback: Authorization header
	if refreshToken == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken := strings.TrimPrefix(authHeader, "Bearer ")
			a.handleAccessTokenFallback(c, accessToken, platform)
			return
		}
	}

	// No valid token found
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}

	// Decode & validate refresh token
	claims, err := DecodeToken(refreshToken, a.Config.JWTSecret)
	if err != nil || claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Generate new tokens
	sub := claims.Sub
	userInfo := map[string]string{
		"sub":     sub,
		"name":    claims.Name,
		"email":   claims.Email,
		"picture": claims.Picture,
	}

	accessToken, newRefreshToken, issuedAt, err := GenerateTokens(sub, time.Duration(a.Config.JWTExpirationSeconds), a.Config.JWTSecret, time.Duration(a.Config.RefreshExpiration), userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	if platform == "web" {
		c.SetCookie(a.Config.CookieName, accessToken, a.Config.JWTExpirationSeconds, "/", "", false, true)
		c.SetCookie(a.Config.RefreshCookieName, newRefreshToken, a.Config.RefreshExpiration, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"issuedAt":  issuedAt,
			"expiresAt": int(issuedAt) + a.Config.JWTExpirationSeconds,
		})
		return
	}

	// Native response
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}

func (a *authRepository) SessionHandler(c *gin.Context) {
	cookieHeader := c.GetHeader("Cookie")
	if cookieHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	cookies := parseCookieHeader(cookieHeader)

	tokenData, ok := cookies[a.Config.CookieName]
	if !ok || tokenData["value"] == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	token := tokenData["value"]

	claims, err := DecodeToken(token, a.Config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var cookieExpiration int64
	if maxAgeStr, ok := cookies[a.Config.CookieName]["maxAge"]; ok {
		maxAge, err := strconv.Atoi(maxAgeStr)
		if err == nil {
			issuedAt := claims.IssuedAt.Time.Unix()
			cookieExpiration = issuedAt + int64(maxAge)
		}
	}

	// Build response from token claims + expiration
	response := gin.H{
		"sub":     claims.Sub,
		"name":    claims.Name,
		"email":   claims.Email,
		"picture": claims.Picture,
		"iat":     claims.IssuedAt.Time.Unix(),
		"exp":     claims.ExpiresAt.Time.Unix(),
	}

	if cookieExpiration > 0 {
		response["cookieExpiration"] = cookieExpiration
	}

	c.JSON(http.StatusOK, response)
}

func (a *authRepository) GoogleLoginHandler(c *gin.Context) {
	redirectURI := c.Query("redirect_uri")
	clientID := c.Query("client_id")
	scope := c.DefaultQuery("scope", "identity")
	state := c.Query("state")

	if redirectURI == "" || clientID != "google" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	platform := "web"
	if redirectURI != a.Config.WebRedirect {
		platform = "mobile"
	}

	googleURL := "https://accounts.google.com/o/oauth2/auth?" + url.Values{
		"client_id":     {a.Config.ClientID},
		"redirect_uri":  {a.Config.ServerURL + "/auth/callback"},
		"response_type": {"code"},
		"scope":         {scope},
		"state":         {platform + "|" + state},
		"prompt":        {"select_account"},
	}.Encode()

	c.Redirect(http.StatusFound, googleURL)
}

func (a *authRepository) AuthGoogleHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	// Step 1: Exchange code for access token
	tokenURL := "https://accounts.google.com/o/oauth2/token"
	data := url.Values{
		"code":          {code},
		"client_id":     {a.Config.ClientID},
		"client_secret": {a.Config.ClientSecret},
		"redirect_uri":  {a.Config.RedirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil || resp.StatusCode >= 400 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get access token"})
		return
	}
	defer resp.Body.Close()

	var tokenData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing access token"})
		return
	}

	// Step 2: Fetch user info
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil || userResp.StatusCode >= 400 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer userResp.Body.Close()

	var userInfo map[string]any
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	c.JSON(http.StatusOK, tokenData)
}

func (a *authRepository) GoogleAuthTokenHandler(c *gin.Context) {
	code := c.PostForm("code")
	fmt.Println("CODE: ", code)
	platform := c.DefaultPostForm("platform", "native")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	// Exchange code for token
	data := url.Values{
		"code":          {code},
		"client_id":     {a.Config.ClientID},
		"client_secret": {a.Config.ClientSecret},
		"redirect_uri":  {a.Config.RedirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil || resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("OAuth error response:", string(bodyBytes))

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "OAuth validation error",
			"message": "Failed to exchange authorization code",
		})
		return
	}
	defer resp.Body.Close()

	var tokenData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode token response"})
		return
	}

	if errMsg, exists := tokenData["error"]; exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             errMsg,
			"error_description": tokenData["error_description"],
			"message":           "OAuth validation error - please ensure the app complies with Google's OAuth 2.0 policy",
		})
		return
	}

	idToken, ok := tokenData["id_token"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required id_token"})
		return
	}

	claims, err := verifyIdToken(idToken, a.Config.ClientID)

	sub := claims.Claims["sub"].(string)
	email := claims.Claims["email"].(string)
	picture := claims.Claims["picture"].(string)
	emailVerified := claims.Claims["email_verified"].(bool)

	if err != nil || sub == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing sub in ID token"})
		return
	}

	userInfo := map[string]string{
		"sub":           sub,
		"email":         email,
		"picture":       picture,
		"emailVerified": strconv.FormatBool(emailVerified),
	}

	accessToken, refreshToken, issuedAt, err := GenerateTokens(sub, time.Duration(a.Config.JWTExpirationSeconds), a.Config.JWTSecret, time.Duration(a.Config.RefreshExpiration), userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate internal tokens"})
		return
	}

	_, err = a.Repo.CreateUser(*a.Ctx, repository.CreateUserParams{
		Sub:           sub,
		Email:         email,
		Picture:       picture,
		EmailVerified: emailVerified,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if platform == "web" {
		// Set cookies
		c.SetCookie(a.Config.CookieName, accessToken, a.Config.JWTExpirationSeconds, "/", "", false, true)
		c.SetCookie(a.Config.RefreshCookieName, refreshToken, a.Config.RefreshExpiration, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"issuedAt":  issuedAt,
			"expiresAt": int(issuedAt) + a.Config.JWTExpirationSeconds,
		})
		return
	}

	// Native (mobile) response
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (a *authRepository) LoginGoogleHandler(c *gin.Context) {
	internalClient := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	requestedScope := c.DefaultQuery("scope", "identity")
	stateParam := c.Query("state")

	if a.Config.ClientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Missing GOOGLE_CLIENT_ID environment variable",
		})
		return
	}

	var platform string
	switch {
	case strings.HasPrefix(redirectURI, "exp://"):
		platform = "mobile"
	case redirectURI == a.Config.WebRedirect:
		platform = "web"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect_uri"})
		return
	}

	if internalClient != "google" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client"})
		return
	}

	if stateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// Construct state for round-trip
	state := fmt.Sprintf("%s|%s", platform, stateParam)

	query := url.Values{
		"client_id":     {a.Config.ClientID},
		"redirect_uri":  {a.Config.ServerURL + "/api/v1/auth/callback"},
		"response_type": {"code"},
		"scope":         {requestedScope},
		"state":         {state},
		"prompt":        {"select_account"},
	}.Encode()

	googleURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?%s", query)
	c.Redirect(http.StatusFound, googleURL)
}

func (a *authRepository) GoogleCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	platform := "web"
	originalState := ""
	if split := strings.SplitN(state, "|", 2); len(split) == 2 {
		platform, originalState = split[0], split[1]
	}

	var redirect string
	if platform == "web" {
		redirect = a.Config.WebRedirect
	} else {
		redirect = a.Config.AppScheme
	}
	redirect += "?code=" + url.QueryEscape(code) + "&state=" + url.QueryEscape(originalState)

	c.Redirect(http.StatusFound, redirect)
}

func (a *authRepository) handleAccessTokenFallback(c *gin.Context, accessToken string, platform string) {
	claims, err := DecodeToken(accessToken, a.Config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	issuedAt := time.Now().Unix()
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Unix(issuedAt, 0))
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(int(issuedAt)+a.Config.JWTExpirationSeconds), 0))

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reissue access token"})
		return
	}

	if platform == "web" {
		c.SetCookie(a.Config.CookieName, newAccessToken, a.Config.JWTExpirationSeconds, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"warning": "Access token used as fallback",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": newAccessToken,
		"warning":     "Access token fallback",
	})
}
