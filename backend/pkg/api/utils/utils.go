package utils

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NullFloat64ToPtr(n sql.NullFloat64) *float64 {
	if n.Valid {
		return &n.Float64
	}
	return nil
}

func BadRequest(c *gin.Context, s string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": s})
	return
}

func InternalServerError(c *gin.Context, s string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": s})
	return
}
