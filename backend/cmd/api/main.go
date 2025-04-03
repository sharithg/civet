package main

// https://github.com/LAA-Software-Engineering/golang-rest-api-template/blob/main/cmd/server/main.go

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	db := database.NewDatabase()
	repo := repository.New(db)
	storage := storage.NewStorage()
	openai := genai.NewOpenAiClient()

	ctx := context.Background()
	logger, _ := zap.NewProduction()

	gin.SetMode(gin.DebugMode)

	r := api.NewRouter(logger, repo, db, storage, openai, &ctx)

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
