package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sharithg/civet/internal/config"
	"github.com/sharithg/civet/internal/repository"
)

func CheckAuth(ctx *context.Context, r *repository.Queries, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenString string

		// Read platform header to determine auth strategy
		platform := c.GetHeader("Platform")

		if strings.ToLower(platform) == "web" {
			// Get token from cookie
			cookie, err := c.Cookie(config.CookieName)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "auth_token cookie missing"})
				c.Abort()
				return
			}
			tokenString = cookie
		} else {
			// Get token from Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
				c.Abort()
				return
			}

			authToken := strings.Split(authHeader, " ")
			if len(authToken) != 2 || authToken[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				c.Abort()
				return
			}
			tokenString = authToken[1]
		}
		// Parse JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		sub := claims["sub"].(string)

		user, err := r.GetUserBySub(*ctx, sub)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}
