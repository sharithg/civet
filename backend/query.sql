-- name: CreateNewOuting :one
INSERT INTO outings (name)
VALUES ($1)
RETURNING id;

-- name: GetReceiptsForOuting :many
SELECT r.restaurant,
    COUNT(oi.id) AS order_count,
    r.total,
    r.id
FROM receipts r
    JOIN order_items oi ON r.id = oi.receipt_id
    JOIN receipt_images ri ON r.receipt_image_id = ri.id
WHERE ri.outing_id = $1
GROUP BY r.id;

-- name: GetOutings :many
SELECT o.id,
    o.name,
    o.created_at,
    COUNT(ri.id) AS total_receipts
FROM outings o
    LEFT JOIN receipt_images ri ON o.id = ri.outing_id
GROUP BY o.id;

-- name: InsertReceiptImage :one
INSERT INTO receipt_images (
        hash,
        bucket,
        key,
        raw_text,
        file_name,
        outing_id
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: InsertReceipt :one
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
RETURNING id;

-- name: InsertOrderItem :exec
INSERT INTO order_items (receipt_id, name, price, quantity)
VALUES ($1, $2, $3, $4);

-- name: InsertOtherFee :exec
INSERT INTO other_fees (receipt_id, name, price)
VALUES ($1, $2, $3);

-- name: GetReceiptByHash :one
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
LIMIT 1;