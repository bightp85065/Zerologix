-- DROP TABLE IF EXISTS
DROP TABLE IF EXISTS order;

-- CREATE TABLE
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id of the order',
    product INT DEFAULT 0 COMMENT 'product ID',
    quantity INT DEFAULT 0 COMMENT 'quantity of the order',
    creator VARCHAR(255) DEFAULT '' COMMENT 'creator of the order',
    action INT DEFAULT 0 COMMENT 'action of the order (1: buy, 2: sell)',
    price_type INT DEFAULT 0 COMMENT 'type of price (1: market, 2: limit)',
    price INT DEFAULT 0 COMMENT 'price of the order',
    status INT DEFAULT 0 COMMENT 'status of the order (1: pending, 2: failure, 3: cancelled, 4: completed)',
    db_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'creation time of the order',
    db_modify_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modification time of the order'
);