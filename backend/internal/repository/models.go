// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	ReceiptID uuid.UUID `json:"receipt_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int32     `json:"quantity"`
}

type OtherFee struct {
	ID        uuid.UUID `json:"id"`
	ReceiptID uuid.UUID `json:"receipt_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
}

type Outing struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type Receipt struct {
	ID                uuid.UUID `json:"id"`
	ReceiptImageID    uuid.UUID `json:"receipt_image_id"`
	Restaurant        string    `json:"restaurant"`
	Address           string    `json:"address"`
	Opened            time.Time `json:"opened"`
	OrderNumber       string    `json:"order_number"`
	OrderType         string    `json:"order_type"`
	TableNumber       string    `json:"table_number"`
	Server            string    `json:"server"`
	Subtotal          float64   `json:"subtotal"`
	SalesTax          float64   `json:"sales_tax"`
	Total             float64   `json:"total"`
	PaymentMethod     string    `json:"payment_method"`
	PaymentAmountPaid float64   `json:"payment_amount_paid"`
	PaymentTip        float64   `json:"payment_tip"`
	Copy              string    `json:"copy"`
	CreatedAt         time.Time `json:"created_at"`
}

type ReceiptImage struct {
	ID       uuid.UUID `json:"id"`
	Bucket   string    `json:"bucket"`
	Key      string    `json:"key"`
	RawText  string    `json:"raw_text"`
	FileName string    `json:"file_name"`
	Hash     string    `json:"hash"`
	OutingID uuid.UUID `json:"outing_id"`
}
