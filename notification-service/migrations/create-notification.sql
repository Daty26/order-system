Create table if not exists notifications(
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL, 
    payment_id INT not null, 
    status VARCHAR(20) not Null,
    message TEXT, 
    created_at TIMESTAMP DEFAULT now()
 )
