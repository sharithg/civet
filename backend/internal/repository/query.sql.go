// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createNewOuting = `-- name: CreateNewOuting :one
INSERT INTO outings (name, user_id, status)
VALUES ($1, $2, $3)
RETURNING id
`

type CreateNewOutingParams struct {
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
	Status string    `json:"status"`
}

func (q *Queries) CreateNewOuting(ctx context.Context, arg CreateNewOutingParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createNewOuting, arg.Name, arg.UserID, arg.Status)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (sub, email, picture, email_verified)
VALUES ($1, $2, $3, $4) ON CONFLICT (sub) DO
UPDATE
SET email = EXCLUDED.email,
    picture = EXCLUDED.picture,
    email_verified = EXCLUDED.email_verified,
    updated_at = NOW()
RETURNING id
`

type CreateUserParams struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Sub,
		arg.Email,
		arg.Picture,
		arg.EmailVerified,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getOutings = `-- name: GetOutings :many
SELECT o.id,
    o.name,
    o.created_at,
    COUNT(ri.id) AS total_receipts,
    COUNT(fr.id) AS total_friends,
    o.status
FROM outings o
    LEFT JOIN receipt_images ri ON o.id = ri.outing_id
    LEFT JOIN friends fr on o.id = fr.outing_id
GROUP BY o.id
`

type GetOutingsRow struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	TotalReceipts int64              `json:"total_receipts"`
	TotalFriends  int64              `json:"total_friends"`
	Status        string             `json:"status"`
}

func (q *Queries) GetOutings(ctx context.Context) ([]GetOutingsRow, error) {
	rows, err := q.db.Query(ctx, getOutings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOutingsRow
	for rows.Next() {
		var i GetOutingsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.TotalReceipts,
			&i.TotalFriends,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReceiptByHash = `-- name: GetReceiptByHash :one
SELECT ri.id AS receipt_image_id,
    ri.hash,
    ri.bucket,
    ri.key,
    ri.raw_text,
    ri.file_name,
    ri.outing_id,
    COALESCE(oi.items, '[]') AS items,
    COALESCE(of.fees, '[]') AS fees
FROM receipt_images ri
    JOIN receipts r ON ri.id = r.receipt_image_id
    LEFT JOIN (
        SELECT receipt_id,
            json_agg(oi.*) AS items
        FROM order_items oi
        GROUP BY receipt_id
    ) oi ON r.id = oi.receipt_id
    LEFT JOIN (
        SELECT receipt_id,
            json_agg(of.*) AS fees
        FROM other_fees of
        GROUP BY receipt_id
    ) of ON r.id = of.receipt_id
WHERE ri.hash = $1
LIMIT 1
`

type GetReceiptByHashRow struct {
	ReceiptImageID uuid.UUID `json:"receipt_image_id"`
	Hash           string    `json:"hash"`
	Bucket         string    `json:"bucket"`
	Key            string    `json:"key"`
	RawText        string    `json:"raw_text"`
	FileName       string    `json:"file_name"`
	OutingID       uuid.UUID `json:"outing_id"`
	Items          []byte    `json:"items"`
	Fees           []byte    `json:"fees"`
}

func (q *Queries) GetReceiptByHash(ctx context.Context, hash string) (GetReceiptByHashRow, error) {
	row := q.db.QueryRow(ctx, getReceiptByHash, hash)
	var i GetReceiptByHashRow
	err := row.Scan(
		&i.ReceiptImageID,
		&i.Hash,
		&i.Bucket,
		&i.Key,
		&i.RawText,
		&i.FileName,
		&i.OutingID,
		&i.Items,
		&i.Fees,
	)
	return i, err
}

const getReceiptsForOuting = `-- name: GetReceiptsForOuting :many
SELECT r.restaurant,
    COUNT(oi.id) AS order_count,
    r.total,
    r.id
FROM receipts r
    JOIN order_items oi ON r.id = oi.receipt_id
    JOIN receipt_images ri ON r.receipt_image_id = ri.id
WHERE ri.outing_id = $1
GROUP BY r.id
`

type GetReceiptsForOutingRow struct {
	Restaurant string    `json:"restaurant"`
	OrderCount int64     `json:"order_count"`
	Total      float64   `json:"total"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) GetReceiptsForOuting(ctx context.Context, outingID uuid.UUID) ([]GetReceiptsForOutingRow, error) {
	rows, err := q.db.Query(ctx, getReceiptsForOuting, outingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReceiptsForOutingRow
	for rows.Next() {
		var i GetReceiptsForOutingRow
		if err := rows.Scan(
			&i.Restaurant,
			&i.OrderCount,
			&i.Total,
			&i.ID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserBySub = `-- name: GetUserBySub :one
select id,
    sub,
    email,
    picture,
    created_at,
    updated_at
from users
where sub = $1
`

type GetUserBySubRow struct {
	ID        uuid.UUID          `json:"id"`
	Sub       string             `json:"sub"`
	Email     string             `json:"email"`
	Picture   string             `json:"picture"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetUserBySub(ctx context.Context, sub string) (GetUserBySubRow, error) {
	row := q.db.QueryRow(ctx, getUserBySub, sub)
	var i GetUserBySubRow
	err := row.Scan(
		&i.ID,
		&i.Sub,
		&i.Email,
		&i.Picture,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertOrderItem = `-- name: InsertOrderItem :exec
INSERT INTO order_items (receipt_id, name, price, quantity)
VALUES ($1, $2, $3, $4)
`

type InsertOrderItemParams struct {
	ReceiptID uuid.UUID `json:"receipt_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int32     `json:"quantity"`
}

func (q *Queries) InsertOrderItem(ctx context.Context, arg InsertOrderItemParams) error {
	_, err := q.db.Exec(ctx, insertOrderItem,
		arg.ReceiptID,
		arg.Name,
		arg.Price,
		arg.Quantity,
	)
	return err
}

const insertOtherFee = `-- name: InsertOtherFee :exec
INSERT INTO other_fees (receipt_id, name, price)
VALUES ($1, $2, $3)
`

type InsertOtherFeeParams struct {
	ReceiptID uuid.UUID `json:"receipt_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
}

func (q *Queries) InsertOtherFee(ctx context.Context, arg InsertOtherFeeParams) error {
	_, err := q.db.Exec(ctx, insertOtherFee, arg.ReceiptID, arg.Name, arg.Price)
	return err
}

const insertReceipt = `-- name: InsertReceipt :one
INSERT INTO receipts (
        receipt_image_id,
        restaurant,
        address,
        opened,
        order_number,
        order_type,
        table_number,
        server,
        subtotal,
        sales_tax,
        total,
        copy
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12
    )
RETURNING id
`

type InsertReceiptParams struct {
	ReceiptImageID uuid.UUID `json:"receipt_image_id"`
	Restaurant     string    `json:"restaurant"`
	Address        string    `json:"address"`
	Opened         time.Time `json:"opened"`
	OrderNumber    string    `json:"order_number"`
	OrderType      string    `json:"order_type"`
	TableNumber    string    `json:"table_number"`
	Server         string    `json:"server"`
	Subtotal       float64   `json:"subtotal"`
	SalesTax       float64   `json:"sales_tax"`
	Total          float64   `json:"total"`
	Copy           string    `json:"copy"`
}

func (q *Queries) InsertReceipt(ctx context.Context, arg InsertReceiptParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertReceipt,
		arg.ReceiptImageID,
		arg.Restaurant,
		arg.Address,
		arg.Opened,
		arg.OrderNumber,
		arg.OrderType,
		arg.TableNumber,
		arg.Server,
		arg.Subtotal,
		arg.SalesTax,
		arg.Total,
		arg.Copy,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const insertReceiptImage = `-- name: InsertReceiptImage :one
INSERT INTO receipt_images (
        hash,
        bucket,
        key,
        raw_text,
        file_name,
        outing_id
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

type InsertReceiptImageParams struct {
	Hash     string    `json:"hash"`
	Bucket   string    `json:"bucket"`
	Key      string    `json:"key"`
	RawText  string    `json:"raw_text"`
	FileName string    `json:"file_name"`
	OutingID uuid.UUID `json:"outing_id"`
}

func (q *Queries) InsertReceiptImage(ctx context.Context, arg InsertReceiptImageParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertReceiptImage,
		arg.Hash,
		arg.Bucket,
		arg.Key,
		arg.RawText,
		arg.FileName,
		arg.OutingID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
