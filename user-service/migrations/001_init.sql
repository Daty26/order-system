Create table users(
    id serial primary key,
    username text unique not null,
    email text unique not null,
    password text not null,
    role text default 'USER',
    created_at TIMESTAMPTZ default now()
)