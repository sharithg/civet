CREATE TABLE splits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    friend_id UUID NOT NULL REFERENCES friends(id),
    order_item_id UUID NOT NULL REFERENCES order_items(id),
    receipt_id UUID NOT NULL REFERENCES receipts(id),
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);