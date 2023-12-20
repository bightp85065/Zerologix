-- DROP TABLE IF EXISTS
DROP TABLE IF EXISTS pair;

-- CREATE TABLE
CREATE TABLE pair (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID of the pair',
    buy_order_id INT DEFAULT 0 COMMENT 'ID of the corresponding buy order',
    sell_order_id INT DEFAULT 0 COMMENT 'ID of the corresponding sell order',
    price INT DEFAULT 0 COMMENT 'price of the pair',
    status INT DEFAULT 0 COMMENT 'status of the pair (1: pending, 2: failure, 3: cancelled, 4: completed)',
    match_time TIMESTAMP COMMENT 'time when the pair was matched',
    db_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'creation time of the pair',
    db_modify_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modification time of the pair'
);