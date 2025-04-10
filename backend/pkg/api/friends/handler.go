package friends

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sharithg/civet/internal/config"
	"github.com/sharithg/civet/internal/repository"
)

type friendsRepository struct {
	DB     *pgxpool.Pool
	Ctx    *context.Context
	Config *config.Config
	Repo   *repository.Queries
}

func New(db *pgxpool.Pool, repo *repository.Queries, config *config.Config, ctx *context.Context) *friendsRepository {
	return &friendsRepository{
		DB:     db,
		Ctx:    ctx,
		Config: config,
		Repo:   repo,
	}
}
