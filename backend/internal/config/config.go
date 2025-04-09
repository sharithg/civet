package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// server
	Port string

	// database
	DbURL string

	// storage
	MinioHost      string
	MinioAccessKey string
	MinioSecretKey string

	// auth
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

	// openai
	OpenAIAPIKey string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: Error loading .env file")
	}

	jwtExpiration, _ := strconv.Atoi(getenv("JWT_EXPIRATION_SECONDS", "900")) // 15 * 60
	refreshExpiration, _ := strconv.Atoi(getenv("REFRESH_EXPIRATION_SECONDS", "604800"))

	cfg := &Config{
		// server
		Port: ":8001",

		// database
		DbURL: envOrPanic("DATABASE_URL"),

		// storage
		MinioHost:      envOrPanic("MINIIO_HOST"),
		MinioAccessKey: envOrPanic("MINIIO_ACCESS_KEY_ID"),
		MinioSecretKey: envOrPanic("MINIIO_SECRET_ACCESS_KEY"),

		// auth
		JWTAlgorithm:         "HS256",
		JWTExpirationSeconds: jwtExpiration,
		RefreshExpiration:    refreshExpiration,

		CookieName:        envOrPanic("COOKIE_NAME"),
		RefreshCookieName: "refresh_token",

		ServerURL:    envOrPanic("SERVER_URL"),
		WebRedirect:  envOrPanic("EXPO_WEB_URL"),
		AppScheme:    envOrPanic("EXPO_APP_SCHEME"),
		ClientID:     envOrPanic("GOOGLE_CLIENT_ID"),
		ClientSecret: envOrPanic("GOOGLE_CLIENT_SECRET"),
		RedirectURI:  envOrPanic("GOOGLE_REDIRECT_URI"),
		JWTSecret:    envOrPanic("JWT_SECRET"),

		// openai
		OpenAIAPIKey: envOrPanic("OPENAI_API_KEY"),
	}

	return cfg
}

func envOrPanic(key string) string {
	fileKey := key + "_FILE"
	if filePath := os.Getenv(fileKey); filePath != "" {
		content := readFromFile(filePath)
		return content
	}

	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func readFromFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", path, err)
	}
	return strings.TrimSpace(string(content))
}

func getenv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
