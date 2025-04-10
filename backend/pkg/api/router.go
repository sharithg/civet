package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sharithg/civet/internal/config"
	"github.com/sharithg/civet/internal/genai"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/internal/storage"
	"github.com/sharithg/civet/pkg/api/auth"
	"github.com/sharithg/civet/pkg/api/outing"
	"github.com/sharithg/civet/pkg/api/receipt"
	"github.com/sharithg/civet/pkg/middleware"
	"go.uber.org/zap"
)

type AppContext struct {
	Logger  *zap.Logger
	Repo    *repository.Queries
	DB      *pgxpool.Pool
	Storage *storage.Storage
	OpenAI  genai.OpenAi
	Context *context.Context
	Config  *config.Config
}

func NewRouter(appCtx *AppContext) *gin.Engine {

	authRepository := auth.New(appCtx.DB, appCtx.Repo, appCtx.Storage, appCtx.OpenAI, appCtx.Config, appCtx.Context)
	outingsRepository := outing.New(appCtx.Repo, appCtx.Context)
	receiptRepository := receipt.New(appCtx.Repo, appCtx.DB, appCtx.Storage, appCtx.OpenAI, appCtx.Context)
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
	v1.Use(middleware.CheckAuth(appCtx.Context, appCtx.Repo, appCtx.Config))

	{
		receipts := v1.Group("/receipt")
		{
			receipts.POST("/upload", receiptRepository.ProcessReceipt)
			receipts.GET("/item/:id", receiptRepository.GetReceipt)
			receipts.POST("/split", receiptRepository.SaveSplit)
			receipts.GET("/:receipt_id/friends", receiptRepository.GetFriends)
			receipts.POST("/friends", receiptRepository.CreateFriend)
			receipts.POST("/friends/split", receiptRepository.CreateSplit)
		}

		outings := v1.Group("/outing")
		{
			outings.POST("", outingsRepository.CreateOuting)
			outings.GET("", outingsRepository.GetOutings)
			outings.GET("/:outing_id/receipts", outingsRepository.GetReceipts)
			outings.GET("/:outing_id/friends", outingsRepository.GetFriends)
		}

	}

	return r
}
