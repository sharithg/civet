-- name: CreateNewOuting :one
INSERT INTO outings (name, user_id, status)
VALUES ($1, $2, $3)
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
    o.status,
    COALESCE(f.friends, '[]') AS friends,
    COALESCE(r.total_receipts, 0) AS total_receipts
FROM outings o
    LEFT JOIN LATERAL (
        SELECT json_agg(json_build_object('id', fr.id, 'name', fr.name)) AS friends
        FROM friends fr
        WHERE fr.outing_id = o.id
    ) f ON true
    LEFT JOIN LATERAL (
        SELECT COUNT(DISTINCT ri.id) AS total_receipts
        FROM receipt_images ri
        WHERE ri.outing_id = o.id
    ) r ON true;

;

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

-- name: CreateUser :one
INSERT INTO users (sub, email, picture, email_verified)
VALUES ($1, $2, $3, $4) ON CONFLICT (sub) DO
UPDATE
SET email = EXCLUDED.email,
    picture = EXCLUDED.picture,
    email_verified = EXCLUDED.email_verified,
    updated_at = NOW()
RETURNING id;

-- name: GetUserBySub :one
select id,
    sub,
    email,
    picture,
    created_at,
    updated_at
from users
where sub = $1;

-- name: GetReceipt :one
SELECT r.id,
    r.total,
    r.restaurant,
    r.address,
    r.opened,
    r.order_number,
    r.order_type,
    r.payment_tip,
    r.payment_amount_paid,
    r.table_number,
    r.copy,
    r.server,
    r.sales_tax,
    ri.bucket,
    ri.key,
    COALESCE(oi.items, '[]') AS items,
    COALESCE(of.fees, '[]') AS fees,
    COALESCE(spl.splits, '[]') AS splits
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
    LEFT JOIN (
        SELECT receipt_id,
            json_agg(
                json_build_object(
                    'id',
                    sp.id,
                    'friend_id',
                    sp.friend_id,
                    'order_item_id',
                    sp.order_item_id
                )
            ) AS splits
        FROM splits sp
        GROUP BY receipt_id
    ) spl ON r.id = spl.receipt_id
WHERE r.id = $1
LIMIT 1;

-- name: CreateOrGetFriend :one
with existing_friend as (
    select id
    from friends
    where friends.user_id = $2
        and friends.outing_id = $3
        and $2 is not null
    limit 1
), inserted_friend as (
    insert into friends (name, user_id, outing_id)
    select $1,
        $2,
        $3
    where not exists (
            select 1
            from existing_friend
        )
    returning id
)
select id
from inserted_friend
union all
select id
from existing_friend
limit 1;

-- name: CreateSplit :one
insert into splits (friend_id, order_item_id, receipt_id, quantity)
values ($1, $2, $3, $4)
returning id;

-- name: DeleteSplit :exec
delete from splits
where receipt_id = $1;

-- name: GetFriends :many
select fr.id,
    fr.name
from friends fr
    join outings o on fr.outing_id = o.id
    join receipt_images ri on o.id = ri.outing_id
    join receipts r on r.receipt_image_id = ri.id
where r.id = $1;

-- name: GetOutingForReceipt :one
select o.id
from outings o
    join receipt_images ri on o.id = ri.outing_id
    join receipts r on ri.id = r.receipt_image_id
where r.id = $1;

-- name: GetReceiptImage :one
select ri.bucket,
    ri.key
from receipt_images ri
    join receipts r on ri.id = r.receipt_image_id
where r.id = $1
limit 1;

-- name: GetCachedCloudVisionResponse :one
select response
from cloud_vision_cache
where image_hash = $1
limit 1;

-- name: GetCachedGenAiResponse :one
select response
from genai_cache
where image_hash = $1
limit 1;

-- name: InsertCachedCloudVisionResponse :one
insert into cloud_vision_cache (image_hash, response)
values ($1, $2)
returning id;

-- name: InsertCachedGenAiResponse :one
insert into genai_cache (image_hash, response)
values ($1, $2)
returning id;


-- name: GetFriendsForOuting :many
WITH unique_friends_per_receipt AS (
    SELECT r.id AS receipt_id, COUNT(DISTINCT fr.id) AS friend_count
    FROM receipts r
    JOIN order_items it ON r.id = it.receipt_id
    JOIN splits sp ON it.id = sp.order_item_id
    JOIN friends fr ON sp.friend_id = fr.id
    GROUP BY r.id
)
SELECT
    fr.name,
    (SUM(it.price * sp.quantity))::float AS subtotal,
    (r.sales_tax / uf.friend_count)::float AS tax_portion,
    (SUM(it.price * sp.quantity) + (r.sales_tax / uf.friend_count))::float AS total_owed
FROM receipts r
JOIN order_items it ON r.id = it.receipt_id
JOIN splits sp ON it.id = sp.order_item_id
JOIN friends fr ON sp.friend_id = fr.id
JOIN receipt_images ri on r.receipt_image_id = ri.id
JOIN outings ou on ri.outing_id = ou.id
JOIN unique_friends_per_receipt uf ON r.id = uf.receipt_id
WHERE ou.id = $1
GROUP BY fr.id, fr.name, r.sales_tax, r.id, uf.friend_count;