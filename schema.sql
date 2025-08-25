CREATE TABLE users (
    id        BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email     VARCHAR(255) NOT NULL UNIQUE,
    password  VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price_cents BIGINT UNSIGNED NOT NULL,
    stock       INT UNSIGNED NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE carts (
    id        BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id   BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_cart (user_id)
);

CREATE TABLE cart_items (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    cart_id      BIGINT UNSIGNED NOT NULL,        
    product_id   BIGINT UNSIGNED NOT NULL,          
    name         VARCHAR(255) NOT NULL,             
    quantity     INT UNSIGNED NOT NULL DEFAULT 1,
    price_cents  BIGINT UNSIGNED NOT NULL,          
    total_cents  BIGINT UNSIGNED NOT NULL,          
    created_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_cart_product (cart_id, product_id)
);

CREATE TABLE orders (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL,
    total_cents BIGINT UNSIGNED NOT NULL,
    status      VARCHAR(50) DEFAULT 'pending',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id   BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    quantity   INT UNSIGNED NOT NULL,
    price_cents BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed user
INSERT INTO users (email, password) VALUES 
('user@example.com', '$2a$10$X4hAbPkzgqGs8AFP0MsgUu35M/GCK93GRE9p1ESFhkexQFvh1IRRe'); -- password: secret123

-- Seed product
INSERT INTO products (name, description, price_cents, stock) VALUES
('Golang T‑Shirt', 'Comfortable cotton tee with Go Gopher', 1999, 100),
('Gin Bottle', 'Premium gin, 750ml', 3499, 50);