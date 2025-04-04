package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sharithg/civet/internal/repository"
)

func GetUser(c *gin.Context) (repository.GetUserBySubRow, error) {
	userRaw, exists := c.Get("currentUser")
	if !exists {
		return repository.GetUserBySubRow{}, errors.New("User not authenticated")
	}

	user, ok := userRaw.(repository.GetUserBySubRow)
	if !ok {
		return repository.GetUserBySubRow{}, errors.New("Failed to cast user")
	}

	return user, nil
}
