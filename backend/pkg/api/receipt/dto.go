package receipt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/pkg/api/utils"
)

type OrderItem struct {
	ID        string   `json:"id"`
	ReceiptID string   `json:"receipt_id"`
	Name      string   `json:"name"`
	Price     *float64 `json:"price"`
	Quantity  int32    `json:"quantity"`
}

type OtherFee struct {
	ID        string   `json:"id"`
	ReceiptID string   `json:"receipt_id"`
	Name      string   `json:"name"`
	Price     *float64 `json:"price"`
}

type Split struct {
	ID       string `json:"id"`
	FriendId string `json:"friend_id"`
	ItemId   string `json:"order_item_id"`
}

type ReceiptResponse struct {
	ID                string      `json:"id"`
	Total             *float64    `json:"total"`
	Restaurant        string      `json:"restaurant"`
	Address           string      `json:"address"`
	Opened            time.Time   `json:"opened"`
	OrderNumber       string      `json:"order_number"`
	OrderType         string      `json:"order_type"`
	PaymentTip        *float64    `json:"payment_tip"`
	PaymentAmountPaid *float64    `json:"payment_amount_paid"`
	TableNumber       string      `json:"table_number"`
	Copy              string      `json:"copy"`
	Server            string      `json:"server"`
	SalesTax          *float64    `json:"sales_tax"`
	Items             []OrderItem `json:"items"`
	ImageUrl          string      `json:"image_url"`
	Fees              []OtherFee  `json:"fees"`
	Splits            []Split     `json:"splits"`
}

func toReceiptResponse(dbRow repository.GetReceiptRow, imageUrl string) ReceiptResponse {
	var items []OrderItem
	if err := json.Unmarshal(dbRow.Items, &items); err != nil {
		log.Printf("error decoding items JSON: %v", err)
		items = []OrderItem{}
	}

	var fees []OtherFee
	if err := json.Unmarshal(dbRow.Fees, &fees); err != nil {
		log.Printf("error decoding fees JSON: %v", err)
		fees = []OtherFee{}
	}

	var splits []Split
	if err := json.Unmarshal(dbRow.Splits, &splits); err != nil {
		log.Printf("error decoding fees JSON: %v", err)
		fees = []OtherFee{}
	}

	return ReceiptResponse{
		ID:                dbRow.ID.String(),
		Total:             utils.NullFloat64ToPtr(dbRow.Total),
		Restaurant:        dbRow.Restaurant,
		Address:           dbRow.Address,
		Opened:            dbRow.Opened,
		OrderNumber:       dbRow.OrderNumber,
		OrderType:         dbRow.OrderType,
		PaymentTip:        utils.NullFloat64ToPtr(dbRow.PaymentTip),
		PaymentAmountPaid: utils.NullFloat64ToPtr(dbRow.PaymentAmountPaid),
		TableNumber:       dbRow.TableNumber,
		Copy:              dbRow.Copy,
		Server:            dbRow.Server,
		SalesTax:          utils.NullFloat64ToPtr(dbRow.SalesTax),
		Items:             items,
		Fees:              fees,
		Splits:            splits,
		ImageUrl:          imageUrl,
	}
}

type SplitItem struct {
	ItemId   string `json:"item_id"`
	Friend   string `json:"friends"`
	Quantity int32  `json:"quantity"`
}

type SplitInput struct {
	ReceiptId string      `json:"receipt_id"`
	Items     []SplitItem `json:"items"`
}

type CreateFriendInput struct {
	ReceiptId string  `json:"receipt_id"`
	Name      string  `json:"name"`
	UserID    *string `json:"user_id"`
}

func toCreateFriend(friend CreateFriendInput, outingId uuid.UUID) (repository.CreateOrGetFriendParams, error) {
	var friendUuid *uuid.UUID
	if friend.UserID == nil {
		friendUuid = nil
	} else {
		uuidVal, err := uuid.Parse(*friend.UserID)
		if err != nil {
			return repository.CreateOrGetFriendParams{}, err
		}
		friendUuid = &uuidVal
	}
	return repository.CreateOrGetFriendParams{
		Name:     friend.Name,
		UserID:   friendUuid,
		OutingID: outingId,
	}, nil
}

type CreateSplitItem struct {
	FriendId string `json:"friend_id"`
	ItemId   string `json:"item_id"`
	Quantity int32  `json:"quantity"`
}
type CreateSplitInput struct {
	ReceiptId string            `json:"receipt_id"`
	Items     []CreateSplitItem `json:"items"`
}

func toCreateSplit(split CreateSplitInput) (*[]repository.CreateSplitParams, *uuid.UUID, error) {
	fmt.Printf("SLIT: %v\n", split)
	receiptUuid, err := uuid.Parse(split.ReceiptId)
	if err != nil {
		return nil, nil, err
	}
	var items []repository.CreateSplitParams
	for _, item := range split.Items {
		itemUuid, err := uuid.Parse(item.ItemId)
		if err != nil {
			return nil, nil, err
		}
		friendUuid, err := uuid.Parse(item.FriendId)
		if err != nil {
			return nil, nil, err
		}
		items = append(items, repository.CreateSplitParams{
			FriendID:    friendUuid,
			OrderItemID: itemUuid,
			Quantity:    item.Quantity,
			ReceiptID:   receiptUuid,
		})
	}
	return &items, &receiptUuid, nil
}
