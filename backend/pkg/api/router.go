package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sharithg/civet/internal/genai"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/internal/storage"
	"github.com/sharithg/civet/pkg/api/auth"
	"github.com/sharithg/civet/pkg/api/outing"
	"github.com/sharithg/civet/pkg/api/receipt"
	"github.com/sharithg/civet/pkg/middleware"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, repo *repository.Queries, db *pgxpool.Pool, storage *storage.Storage, openai genai.OpenAi, ctx *context.Context) *gin.Engine {
	config, err := auth.LoadConfig()

	if err != nil {
		log.Fatalln("error loading config")
	}

	receiptRepository := receipt.New(repo, db, storage, openai, ctx)
	authRepository := auth.New(db, storage, openai, config, ctx)
	outingsRepository := outing.New(repo, ctx)

	r := gin.Default()

	r.Use(middleware.Cors())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/receipt/upload", receiptRepository.ProcessReceipt)

		auth := v1.Group("/auth")
		{
			auth.GET("/callback", authRepository.GoogleCallbackHandler)
			auth.GET("/login/google", authRepository.LoginGoogleHandler)
			auth.GET("/google", authRepository.AuthGoogleHandler)
			auth.POST("/token", authRepository.GoogleAuthTokenHandler)
			auth.POST("/refresh", authRepository.RefreshTokenHandler)
			auth.GET("/session", authRepository.SessionHandler)
		}

		outings := v1.Group("/outing")
		{
			outings.POST("", outingsRepository.CreateOuting)
			outings.GET("", outingsRepository.GetOutings)
			outings.GET("/:outing_id/receipts", outingsRepository.GetReceipts)
		}
	}

	return r
}
