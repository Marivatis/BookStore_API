CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    author VARCHAR(255),
    price NUMERIC(10, 2),
    stock INT,
    created_at TIMESTAMP
);