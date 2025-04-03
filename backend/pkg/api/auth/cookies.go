package auth

import (
	"net/http"
	"strings"
)

func BuildCookie(name, value string, maxAge int, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func parseCookieHeader(cookieHeader string) map[string]map[string]string {
	result := map[string]map[string]string{}

	parts := strings.Split(cookieHeader, ";")
	var lastKey string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "=") {
			keyVal := strings.SplitN(part, "=", 2)
			key := keyVal[0]
			val := keyVal[1]
			result[key] = map[string]string{"value": val}
			lastKey = key
		} else if strings.HasPrefix(strings.ToLower(part), "max-age=") && lastKey != "" {
			result[lastKey]["maxAge"] = strings.TrimPrefix(part, "max-age=")
		} else if strings.HasPrefix(strings.ToLower(part), "expires=") && lastKey != "" {
			result[lastKey]["expires"] = strings.TrimPrefix(part, "expires=")
		} else if strings.ToLower(part) == "httponly" && lastKey != "" {
			result[lastKey]["httpOnly"] = "true"
		}
	}

	return result
}
