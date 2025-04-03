package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Sub     string `json:"sub"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Picture string `json:"picture,omitempty"`
	Type    string `json:"type,omitempty"` // "refresh" for refresh token
}

func GenerateTokens(sub string, jwtExp time.Duration, jwtSecret string, refreshExp time.Duration, userInfo map[string]string) (accessToken, refreshToken string, issuedAt int64, err error) {
	now := time.Now()
	issuedAt = now.Unix()
	jti := uuid.New().String()

	claims := CustomClaims{
		Sub:     sub,
		Name:    userInfo["name"],
		Email:   userInfo["email"],
		Picture: userInfo["picture"],
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtExp * time.Second)),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", 0, err
	}

	claims.Type = "refresh"
	claims.RegisteredClaims.ID = jti
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(refreshExp * time.Second))

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	return
}

func DecodeToken(tokenStr string, jwtSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

func verifyIdToken(idToken string, aud string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, aud)

	if err != nil {
		return nil, err
	}
	return payload, nil
}
