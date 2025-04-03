CREATE TABLE outings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE receipt_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bucket VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL,
    raw_text TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    hash TEXT NOT NULL,
    outing_id UUID NOT NULL REFERENCES outings(id)
);

CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_image_id UUID NOT NULL REFERENCES receipt_images(id) ON DELETE CASCADE,
    restaurant VARCHAR(255) NOT NULL,
    address TEXT,
    opened TIMESTAMP,
    order_number VARCHAR(50),
    order_type VARCHAR(50),
    table_number VARCHAR(50),
    server VARCHAR(255),
    subtotal NUMERIC(10, 2) NOT NULL,
    sales_tax NUMERIC(10, 2) NOT NULL,
    total NUMERIC(10, 2) NOT NULL,
    payment_method VARCHAR(50),
    payment_amount_paid NUMERIC(10, 2),
    payment_tip NUMERIC(10, 2),
    copy VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    quantity INT NOT NULL
);

CREATE TABLE other_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);