create table if not exists inventory(
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
)