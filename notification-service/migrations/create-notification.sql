create table IF NOT EXISTS notifications (
  id SERIAL PRIMARY KEY,
  order_id int not null,
  payment_id int not null,
  status varchar(20) not null check (status IN ('PENDING','SENT','FAILED')),
  message TEXT,
  created_at timestamp default now()
);

