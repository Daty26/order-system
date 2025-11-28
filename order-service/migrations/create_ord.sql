Create table orders(
   id serial primary key,
   user_id int not null,
   status VARCHAR(50) not null default 'PENDING',
   created_at TIMESTAMP not null default now()
)

create table order_items(
    id serial primary key,
    order_id int not null references orders(id) on delete cascade,
    product_id int not null,
    quantity int not null check( quantity >0)
)