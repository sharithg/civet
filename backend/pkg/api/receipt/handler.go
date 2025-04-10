package receipt

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sharithg/civet/internal/config"
	"github.com/sharithg/civet/internal/genai"
	"github.com/sharithg/civet/internal/receipt"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/internal/storage"
	"github.com/sharithg/civet/pkg/api/utils"
)

type ReceiptWithDetails struct {
	ImageID  uuid.UUID              `json:"image_id"`
	Hash     string                 `json:"hash"`
	Bucket   string                 `json:"bucket"`
	Key      string                 `json:"key"`
	RawText  string                 `json:"raw_text"`
	FileName string                 `json:"file_name"`
	OutingID uuid.UUID              `json:"outing_id,omitempty"`
	Items    []repository.OrderItem `json:"items"`
	Fees     []repository.OtherFee  `json:"fees"`
}

type receiptRepository struct {
	Repo    *repository.Queries
	Ctx     *context.Context
	Storage *storage.Storage
	Genai   genai.OpenAi
	Db      *pgxpool.Pool
	Config  *config.Config
}

func New(repo *repository.Queries, db *pgxpool.Pool, storage *storage.Storage, genai genai.OpenAi, ctx *context.Context, config *config.Config) *receiptRepository {
	return &receiptRepository{
		Repo:    repo,
		Ctx:     ctx,
		Storage: storage,
		Genai:   genai,
		Db:      db,
		Config:  config,
	}
}

func (r *receiptRepository) SaveReceipt(repo *repository.Queries, hash, bucket, key, text, name string, outingId uuid.UUID, receipt receipt.ParsedReceipt) error {
	tx, err := r.Db.BeginTx(*r.Ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(*r.Ctx)

	qtx := repo.WithTx(tx)

	// 1. Insert into receipt_images
	imageId, err := qtx.InsertReceiptImage(*r.Ctx, repository.InsertReceiptImageParams{
		Hash:     hash,
		Bucket:   bucket,
		Key:      key,
		RawText:  text,
		FileName: name,
		OutingID: outingId,
	})
	if err != nil {
		return fmt.Errorf("insert into receipt_images: %w", err)
	}

	// 2. Insert into receipts
	receiptId, err := qtx.InsertReceipt(*r.Ctx, repository.InsertReceiptParams{
		ReceiptImageID: imageId,
		Restaurant:     receipt.Restaurant,
		Address:        receipt.Address,
		Opened:         receipt.Opened,
		OrderNumber:    receipt.OrderNumber,
		OrderType:      receipt.OrderType,
		TableNumber:    receipt.Table,
		Server:         receipt.Server,
		Subtotal: sql.NullFloat64{
			Float64: receipt.Subtotal,
		},
		SalesTax: sql.NullFloat64{
			Float64: receipt.SalesTax,
		},
		Total: sql.NullFloat64{
			Float64: receipt.Total,
		},
		Copy: receipt.Copy,
	})
	if err != nil {
		return fmt.Errorf("insert into receipts: %w", err)
	}

	// 3. Insert order_items
	for _, item := range receipt.Items {
		err = qtx.InsertOrderItem(*r.Ctx, repository.InsertOrderItemParams{
			ReceiptID: receiptId,
			Name:      item.Name,
			Price: sql.NullFloat64{
				Float64: item.Price,
			},
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			return fmt.Errorf("insert into order_items: %w", err)
		}
	}

	// 4. Insert other_fees
	for _, fee := range receipt.OtherFees {
		err = qtx.InsertOtherFee(*r.Ctx, repository.InsertOtherFeeParams{
			ReceiptID: receiptId,
			Name:      fee.Name,
			Price: sql.NullFloat64{
				Float64: fee.Price,
			},
		})
		if err != nil {
			return fmt.Errorf("insert into other_fees: %w", err)
		}
	}

	// 5. Commit transaction
	if err := tx.Commit(*r.Ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *receiptRepository) GetReceiptByHash(hash string) (*ReceiptWithDetails, error) {
	row, err := r.Repo.GetReceiptByHash(*r.Ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	var items []repository.OrderItem
	var fees []repository.OtherFee

	if err := json.Unmarshal([]byte(row.Items), &items); err != nil {
		return nil, fmt.Errorf("unmarshal items: %w", err)
	}
	if err := json.Unmarshal([]byte(row.Fees), &fees); err != nil {
		return nil, fmt.Errorf("unmarshal fees: %w", err)
	}

	return &ReceiptWithDetails{
		ImageID:  row.ReceiptImageID,
		Hash:     row.Hash,
		Bucket:   row.Bucket,
		Key:      row.Key,
		RawText:  row.RawText,
		FileName: row.FileName,
		OutingID: row.OutingID,
		Items:    items,
		Fees:     fees,
	}, nil
}

func (r *receiptRepository) ProcessReceipt(c *gin.Context) {
	fileHeader, err := c.FormFile("photo.0")

	outingId := uuid.MustParse(c.GetHeader("outingId"))

	if err != nil {
		fmt.Println("Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
		fmt.Println("Error on file type: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("Error on opening file: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "opening file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error on reading file data: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "reading file data"})
		return
	}

	fileInfo, err := receipt.NewExtract(*r.Ctx, *r.Storage, r.Genai, r.Repo, data, fileHeader.Filename, r.Config.CloudVisionCredentials)

	if err != nil {
		fmt.Println("Error on starting extraction: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "starting extraction"})
		return
	}

	existing, err := r.GetReceiptByHash(fileInfo.ImageHash)

	if err != nil {
		fmt.Println("Error on getting existing receipt image: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing receipt image"})
		return
	}

	if existing != nil {
		c.JSON(http.StatusOK, gin.H{"hash": fileInfo.ImageHash, "existing": true})
		return
	}

	model, text, bucket, key, err := fileInfo.Run(*r.Ctx)

	if err != nil {
		fmt.Println("Error on running model: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload receipt image"})
		return
	}

	if err = r.SaveReceipt(r.Repo, fileInfo.ImageHash, bucket, key, text, fileHeader.Filename, outingId, model); err != nil {
		fmt.Println("Error on saving receipt: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save receipt image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hash": fileInfo.ImageHash, "existing": false})
}

func (r *receiptRepository) GetReceipt(c *gin.Context) {
	receiptId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receipt id"})
		return
	}

	receipt, err := r.Repo.GetReceipt(*r.Ctx, receiptId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fetching receipt"})
		return
	}

	url, err := r.Storage.GetObjectUrl(*r.Ctx, receipt.Bucket, receipt.Key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "getting object url"})
		return
	}

	c.JSON(http.StatusOK, toReceiptResponse(receipt, url))
}

func (r *receiptRepository) SaveSplit(c *gin.Context) {

	var body SplitInput
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "invalid request")
		return
	}

	for _, item := range body.Items {
		itemUuid, err := uuid.Parse(item.ItemId)
		if err != nil {
			utils.BadRequest(c, "invalid item id")
			return
		}
		friendUuid, err := uuid.Parse(item.Friend)
		if err != nil {
			utils.BadRequest(c, "invalid friend id")
			return
		}
		r.Repo.CreateSplit(*r.Ctx, repository.CreateSplitParams{
			FriendID:    friendUuid,
			OrderItemID: itemUuid,
			Quantity:    item.Quantity,
		})
	}
}

func (r *receiptRepository) GetFriends(c *gin.Context) {
	receiptIdStr := c.Param("receipt_id")
	receiptId, err := uuid.Parse(receiptIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outing ID"})
		return
	}

	friends, err := r.Repo.GetFriends(*r.Ctx, receiptId)
	if err != nil {
		fmt.Println("ERR: ", err)
		utils.InternalServerError(c, "failed to fetch friends")
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (r *receiptRepository) CreateFriend(c *gin.Context) {

	var body CreateFriendInput
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("ERR: ", err)
		utils.BadRequest(c, "invalid request")
		return
	}

	receiptUuid, err := uuid.Parse(body.ReceiptId)
	if err != nil {
		utils.BadRequest(c, "invalid receipt id")
		return
	}

	outing, err := r.Repo.GetOutingForReceipt(*r.Ctx, receiptUuid)
	if err != nil {
		utils.BadRequest(c, "invalid outing id")
		return
	}

	createFriend, err := toCreateFriend(body, outing)
	if err != nil {
		fmt.Println("ERR: ", err)
		utils.BadRequest(c, "invalid request")
		return
	}

	friendId, err := r.Repo.CreateOrGetFriend(*r.Ctx, createFriend)
	if err != nil {
		fmt.Println("ERR: ", err)
		utils.InternalServerError(c, "failed to create friend")
		return
	}

	c.JSON(http.StatusOK, gin.H{"friendId": friendId})
}

func (r *receiptRepository) CreateSplit(c *gin.Context) {

	var body CreateSplitInput
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("ERR b: ", err)
		utils.BadRequest(c, "invalid request")
		return
	}

	createSplit, receiptId, err := toCreateSplit(body)
	if err != nil {
		utils.BadRequest(c, "invalid request")
		return
	}

	if err = r.Repo.DeleteSplit(*r.Ctx, *receiptId); err != nil {
		utils.BadRequest(c, "error deleting split")
		return
	}

	for _, split := range *createSplit {
		_, err := r.Repo.CreateSplit(*r.Ctx, split)
		if err != nil {
			utils.InternalServerError(c, "failed to create split")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
