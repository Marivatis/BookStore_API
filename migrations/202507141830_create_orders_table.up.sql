CREATE TYPE order_status AS ENUM (
    'created',
    'accepted',
    'pending',
    'paid',
    'shipped',
    'delivered',
    'cancelled'
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    status order_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP
);