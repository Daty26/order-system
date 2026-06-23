ALTER TABLE orders
ADD COLUMN total_amount_cents BIGINT NOT NULL DEFAULT 0
CHECK (total_amount_cents >= 0);

ALTER TABLE order_items
ADD COLUMN unit_price_cents BIGINT NOT NULL DEFAULT 0
CHECK (unit_price_cents >= 0);