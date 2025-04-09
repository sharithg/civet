package outing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/pkg/api/auth"
	"github.com/sharithg/civet/pkg/api/utils"
)

type Repository struct {
	Repo *repository.Queries
	Ctx  *context.Context
}

func New(repo *repository.Queries, ctx *context.Context) *Repository {
	return &Repository{Repo: repo, Ctx: ctx}
}

type CreateOutingRequest struct {
	Name string `json:"name" binding:"required"`
}

func (r *Repository) CreateOuting(c *gin.Context) {

	user, err := auth.GetUser(c)

	if err != nil {
		fmt.Println("err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "getting user"})
		return
	}

	var body CreateOutingRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, err := r.Repo.CreateNewOuting(*r.Ctx, repository.CreateNewOutingParams{
		Name:   body.Name,
		UserID: user.ID,
		Status: "active",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (r *Repository) GetOutings(c *gin.Context) {
	outings, err := r.Repo.GetOutings(*r.Ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch outings"})
		return
	}

	outingsResp, err := toOutingsResponse(outings)

	if err != nil {
		utils.BadRequest(c, "error reading outings")
		return
	}

	c.JSON(http.StatusOK, outingsResp)
}

func (r *Repository) GetReceipts(c *gin.Context) {
	outingIDStr := c.Param("outing_id")
	outingID, err := uuid.Parse(outingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outing ID"})
		return
	}

	receipts, err := r.Repo.GetReceiptsForOuting(*r.Ctx, outingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch receipts"})
		return
	}
	c.JSON(http.StatusOK, toOutingReceiptsResponse(receipts))
}
