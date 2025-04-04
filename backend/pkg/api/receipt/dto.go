package receipt

import (
	"encoding/json"
	"log"
	"time"

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
	Fees              []OtherFee  `json:"fees"`
}

func toReceiptResponse(dbRow repository.GetReceiptRow) ReceiptResponse {
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
	}
}
