CREATE UNIQUE INDEX IF NOT EXISTS payments_order_id_unique
ON payments(order_id);
