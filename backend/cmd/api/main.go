package main

// https://github.com/LAA-Software-Engineering/golang-rest-api-template/blob/main/cmd/server/main.go

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sharithg/civet/internal/config"
	"github.com/sharithg/civet/internal/genai"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/internal/storage"
	"github.com/sharithg/civet/pkg/api"
	"github.com/sharithg/civet/pkg/database"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.LoadConfig()

	db := database.NewDatabase(config)
	repo := repository.New(db)
	storage := storage.NewStorage(config)
	openai := genai.NewOpenAiClient(config)

	ctx := context.Background()
	logger, _ := zap.NewProduction()

	gin.SetMode(gin.DebugMode)

	appCtx := api.AppContext{
		Logger:  logger,
		Repo:    repo,
		DB:      db,
		Storage: storage,
		OpenAI:  openai,
		Context: &ctx,
		Config:  config,
	}

	r := api.NewRouter(&appCtx)

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
