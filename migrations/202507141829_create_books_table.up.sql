CREATE TABLE books (
    product_id INT PRIMARY KEY,
    author VARCHAR(255),
    isbn VARCHAR(20),
    CONSTRAINT fk_product_book
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);