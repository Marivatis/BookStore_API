CREATE TABLE magazines (
    product_id INT PRIMARY KEY,
    issue_number INT,
    publication_date DATE,
    CONSTRAINT fk_product_magazine
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);