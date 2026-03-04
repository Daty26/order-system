CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    user_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
