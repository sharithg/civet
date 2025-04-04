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

	authRepository := auth.New(db, repo, storage, openai, config, ctx)

	r := gin.Default()

	r.Use(middleware.Cors())

	// Auth routes
	authRoutes := r.Group("/api/v1/auth")
	{
		authRoutes.GET("/callback", authRepository.GoogleCallbackHandler)
		authRoutes.GET("/login/google", authRepository.LoginGoogleHandler)
		authRoutes.GET("/google", authRepository.AuthGoogleHandler)
		authRoutes.POST("/token", authRepository.GoogleAuthTokenHandler)
		authRoutes.POST("/refresh", authRepository.RefreshTokenHandler)
		authRoutes.GET("/session", authRepository.SessionHandler)
	}

	v1 := r.Group("/api/v1")
	v1.Use(middleware.CheckAuth(ctx, repo))

	outingsRepository := outing.New(repo, ctx)
	receiptRepository := receipt.New(repo, db, storage, openai, ctx)

	{
		v1.POST("/receipt/upload", receiptRepository.ProcessReceipt)

		outings := v1.Group("/outing")
		{
			outings.POST("", outingsRepository.CreateOuting)
			outings.GET("", outingsRepository.GetOutings)
			outings.GET("/:outing_id/receipts", outingsRepository.GetReceipts)
		}
	}

	return r
}
