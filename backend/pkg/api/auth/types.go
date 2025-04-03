package auth

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	JWTAlgorithm         string
	JWTExpirationSeconds int
	RefreshExpiration    int

	CookieName        string
	RefreshCookieName string

	ServerURL    string
	WebRedirect  string
	AppScheme    string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	JWTSecret    string
}

func LoadConfig() (*Config, error) {
	jwtExpiration, _ := strconv.Atoi(getenv("JWT_EXPIRATION_SECONDS", "900")) // 15 * 60
	refreshExpiration, _ := strconv.Atoi(getenv("REFRESH_EXPIRATION_SECONDS", "604800"))

	cfg := &Config{
		JWTAlgorithm:         "HS256",
		JWTExpirationSeconds: jwtExpiration,
		RefreshExpiration:    refreshExpiration,

		CookieName:        "auth_token",
		RefreshCookieName: "refresh_token",

		ServerURL:    getenv("SERVER_URL", "http://localhost:8001"),
		WebRedirect:  getenv("EXPO_WEB_URL", "http://localhost:8081"),
		AppScheme:    getenv("EXPO_APP_SCHEME", "exp://10.0.0.63:8081"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}

	// Validate required fields
	missing := []string{}
	if cfg.ClientID == "" {
		missing = append(missing, "GOOGLE_CLIENT_ID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "GOOGLE_CLIENT_SECRET")
	}
	if cfg.RedirectURI == "" {
		missing = append(missing, "GOOGLE_REDIRECT_URI")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func getenv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
