CREATE TABLE order_items (
    order_id SERIAL NOT NULL,
    book_id SERIAL NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price NUMERIC(10,2) NOT NULL,
    PRIMARY KEY (order_id, book_id),
    CONSTRAINT fk_order
         FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_books
         FOREIGN KEY (book_id) REFERENCES books(id)
);