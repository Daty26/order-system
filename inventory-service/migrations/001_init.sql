create table if not exists inventory(
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
)